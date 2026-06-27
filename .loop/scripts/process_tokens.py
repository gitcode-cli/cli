#!/usr/bin/env python3
"""Process claude -p stream-json output: extract readable text, token data, and inject into delivery files."""
import json, re, sys, os, subprocess
from datetime import datetime, timedelta, timezone

def main():
    jsonl_file = sys.argv[1]
    log_file = sys.argv[2]
    token_file = sys.argv[3]
    deliveries_dir = sys.argv[4]

    with open(jsonl_file) as f:
        lines = f.readlines()

    # Step 1: Extract readable text → log
    with open(log_file, 'a') as log:
        for line in lines:
            line = line.strip()
            if not line.startswith('{'):
                continue
            try:
                d = json.loads(line)
            except json.JSONDecodeError:
                continue
            if d.get('type') == 'assistant':
                msg = d.get('message', {})
                for c in msg.get('content', []):
                    if c.get('type') == 'text':
                        log.write(c['text'] + '\n')

    # Step 2: Find result and extract token data
    result = None
    for line in lines:
        line = line.strip()
        if not line.startswith('{'):
            continue
        try:
            d = json.loads(line)
        except json.JSONDecodeError:
            continue
        if d.get('type') == 'result':
            result = d

    if not result:
        with open(log_file, 'a') as log:
            log.write('WARNING: no result object found in stream-json output\n')
        return

    u = result.get('usage', {})
    mu = result.get('modelUsage', {})
    total_in = u.get('input_tokens', 0)
    total_out = u.get('output_tokens', 0)
    cache_read = u.get('cache_read_input_tokens', 0)
    cache_create = u.get('cache_creation_input_tokens', 0)
    dur_ms = result.get('duration_ms', 0)
    turns = result.get('num_turns', 0)

    # Calculate DeepSeek cost
    cost_rmb = 0.0
    for m, v in mu.items():
        mi = v.get('inputTokens', 0)
        mc = v.get('cacheReadInputTokens', 0)
        mo = v.get('outputTokens', 0)
        cost_rmb += (mi / 1_000_000) * 3.0
        cost_rmb += (mc / 1_000_000) * 0.025
        cost_rmb += (mo / 1_000_000) * 6.0
    cost_rmb = round(cost_rmb, 4)
    cost_usd = round(cost_rmb / 7.2, 4)

    # Log summary
    with open(log_file, 'a') as log:
        log.write('---\n')
        log.write(f'TOKENS:  in={total_in}  out={total_out}  total={total_in+total_out}  cache_read={cache_read}  cache_create={cache_create}  cost=¥{cost_rmb}  duration={dur_ms}ms  turns={turns}\n')
        for m, v in mu.items():
            log.write(f'  {m}: in={v.get("inputTokens",0)} out={v.get("outputTokens",0)} cache_read={v.get("cacheReadInputTokens",0)} cost=¥{cost_rmb}\n')

    # Write token JSON
    token_data = {
        'input_tokens': total_in,
        'cache_miss_tokens': total_in,
        'cache_read_input_tokens': cache_read,
        'cache_creation_input_tokens': cache_create,
        'output_tokens': total_out,
        'total_tokens': total_in + total_out,
        'cost_rmb': cost_rmb,
        'cost_usd': cost_usd,
        'duration_ms': dur_ms,
        'num_turns': turns,
        'pricing': 'DeepSeek: ¥3/M cache-miss, ¥0.025/M cache-hit, ¥6/M output',
    }
    with open(token_file, 'w') as f:
        json.dump(token_data, f, indent=2)

    # Step 3: Extract ISSUE_NUM
    with open(jsonl_file) as f:
        content = f.read()
    m = re.search(r'ISSUE_NUM=(\d+)', content)
    if not m:
        # Fallback to log
        try:
            with open(log_file) as f:
                content = f.read()
            m = re.search(r'ISSUE_NUM=(\d+)', content)
        except:
            pass
    if not m:
        with open(log_file, 'a') as log:
            log.write('no issue num found\n')
        return

    issue_num = m.group(1)

    # Step 4: Inject token data into delivery files
    delivery_file = os.path.join(deliveries_dir, f'issue-{issue_num}.md')
    readme_file = os.path.join(deliveries_dir, 'README.md')

    # Build token block
    in_k = total_in / 1000
    out_k = total_out / 1000
    total_k = (total_in + total_out) / 1000
    cache_read_k = cache_read / 1000
    dur = f'{dur_ms/1000:.0f}s'

    token_block = f"""
## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | {total_in:,} ({in_k:.0f}k) |
| 输出 tokens | {total_out:,} ({out_k:.0f}k) |
| 缓存命中 | {cache_read:,} ({cache_read_k:.0f}k) |
| 缓存写入 | {cache_create:,} |
| 总计 tokens | {total_in+total_out:,} ({total_k:.0f}k) |
| 成本 (DeepSeek) | ¥{cost_rmb} (~${cost_usd}) |
| 耗时 | {dur} |
| 轮次 | {turns} |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
"""

    # Write to delivery file (create if missing — subprocess may have lost it in worktree)
    if not os.path.exists(delivery_file):
        with open(delivery_file, 'w') as f:
            f.write(f'# Issue #{issue_num} — Delivery Record\n\n')
    with open(delivery_file, 'a') as f:
        f.write(token_block)
    with open(log_file, 'a') as log:
        log.write(f'token data written to {delivery_file}\n')

    # Update README row — total includes cache
    total_with_cache = total_in + total_out + cache_read
    if total_with_cache >= 1_000_000:
        total_str = f'{total_with_cache/1_000_000:.1f}M'
    else:
        total_str = f'{total_with_cache/1000:.0f}k'
    cache_str = f'{cache_read/1_000_000:.1f}M' if cache_read >= 1_000_000 else f'{cache_read/1000:.0f}k'
    token_str = f'{total_str}({cache_str} cache)'
    cost_str = f'¥{cost_rmb}'

    # Look up git stats for this issue
    pr_str = '—'
    change_str = '—'
    time_str = '—'
    since_date = (datetime.now(timezone.utc) - timedelta(days=60)).strftime('%Y-%m-%d')
    try:
        sha_result = subprocess.run(
            ['git', 'log', 'origin/main', '--merges', '--format=%H|%ci|%s', f'--since={since_date}'],
            capture_output=True, text=True, cwd='/home/wpf/claude-code/vibe-coding/cli'
        )
        for line in sha_result.stdout.strip().split('\n'):
            if not line or f'issue-{issue_num}' not in line.lower():
                continue
            parts = line.split('|', 2)
            sha = parts[0]
            time_str = parts[1][:16] if len(parts) > 1 else '—'
            msg = parts[2] if len(parts) > 2 else ''
            pr_match = re.search(r'!(\d+)', msg)
            if pr_match:
                pr_str = f'#{pr_match.group(1)}'
            diff_result = subprocess.run(
                ['git', 'diff', f'{sha}~1..{sha}', '--stat'],
                capture_output=True, text=True, cwd='/home/wpf/claude-code/vibe-coding/cli'
            )
            m = re.search(r'(\d+)\s+files?\s+changed(?:,\s+(\d+)\s+insertions?\(\+\))?(?:,\s+(\d+)\s+deletions?\(-\))?', diff_result.stdout)
            if m:
                adds = int(m.group(2)) if m.group(2) else 0
                dels = int(m.group(3)) if m.group(3) else 0
                change_str = f'+{adds}/-{dels}'
            break
    except Exception:
        pass

    if os.path.exists(readme_file):
        with open(readme_file) as f:
            lines = f.readlines()
        target = f'| [#{issue_num}](issue-{issue_num}.md) |'
        matches = []
        for i, line in enumerate(lines):
            if target in line:
                # Check it's a table row (starts with |), not a reference link
                stripped = line.strip()
                if stripped.startswith('|') and stripped.endswith('|'):
                    matches.append(i)
        if matches:
            # Keep the last match (most recent), dedup by removing earlier ones
            keep_idx = matches[-1]
            stale_idxs = matches[:-1]
            # Remove stale rows (iterate reversed to preserve indices)
            for idx in reversed(stale_idxs):
                del lines[idx]
                # Adjust keep_idx if we deleted a row before it
                if idx < keep_idx:
                    keep_idx -= 1
            # Update token/cost/time/change columns on the kept row
            cols = [c.strip() for c in lines[keep_idx].split('|') if c.strip()]
            if len(cols) >= 11:
                cols[-5] = token_str     # tokens
                cols[-4] = change_str    # change
                cols[-2] = cost_str      # cost
                cols[-1] = time_str      # time
                if pr_str != '—': cols[-8] = pr_str  # PR
            lines[keep_idx] = '| ' + ' | '.join(cols) + ' |\n'
            with open(readme_file, 'w') as f:
                f.writelines(lines)
            with open(log_file, 'a') as log:
                if stale_idxs:
                    log.write(f'README dedup: removed {len(stale_idxs)} stale row(s) for #{issue_num}, updated token={token_str}\n')
                else:
                    log.write(f'README row updated with token={token_str}\n')
        else:
            # Row not found — insert new row before ## 统计
            new_row = f'| [#{issue_num}](issue-{issue_num}.md) | — | merged | {pr_str} | — | — | — | {change_str} | {token_str} | {cost_str} | {time_str} |\n'
            for j, l in enumerate(lines):
                if l.startswith('## 统计'):
                    lines.insert(j, new_row)
                    break
            with open(readme_file, 'w') as f:
                f.writelines(lines)
            with open(log_file, 'a') as log:
                log.write(f'README new row added for #{issue_num}\n')

    # Refresh stats
    subprocess.run(['bash', 'scripts/count-deliveries.sh'], cwd='/home/wpf/claude-code/vibe-coding/cli')

    print(f'ISSUE_NUM={issue_num} tokens={total_in+total_out} cost=¥{cost_rmb}')


if __name__ == '__main__':
    main()
