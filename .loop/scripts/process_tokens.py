#!/usr/bin/env python3
"""Process claude -p stream-json output: extract readable text, token data, and inject into delivery files."""
import json, re, sys, os, subprocess

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

    # Append to delivery file
    if os.path.exists(delivery_file):
        with open(delivery_file, 'a') as f:
            f.write(token_block)
        with open(log_file, 'a') as log:
            log.write(f'token data appended to {delivery_file}\n')

    # Update README row
    token_str = f'{total_k:.0f}k'
    if os.path.exists(readme_file):
        with open(readme_file) as f:
            lines = f.readlines()
        target = f'| [#{issue_num}](issue-{issue_num}.md) |'
        for i, line in enumerate(lines):
            if target in line:
                cols = [c.strip() for c in line.split('|') if c.strip()]
                # Update token column
                if len(cols) >= 9:
                    cols[-1] = token_str
                lines[i] = '| ' + ' | '.join(cols) + ' |\n'
                break
        with open(readme_file, 'w') as f:
            f.writelines(lines)
        with open(log_file, 'a') as log:
            log.write(f'README row updated with token={token_str}\n')

    # Refresh stats
    subprocess.run(['bash', 'scripts/count-deliveries.sh'], cwd='/home/wpf/claude-code/vibe-coding/cli')

    print(f'ISSUE_NUM={issue_num} tokens={total_in+total_out} cost=¥{cost_rmb}')


if __name__ == '__main__':
    main()
