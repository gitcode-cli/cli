import importlib.util
import json
import sys
import unittest
from pathlib import Path
from unittest import mock

# verify-remote-facts.py has a hyphen in its name so it cannot be imported
# normally; load it via importlib and register it so mock.patch can find it.
_spec = importlib.util.spec_from_file_location(
    "verify_remote_facts",
    Path(__file__).parent.parent / "scripts" / "verify-remote-facts.py",
)
vrf = importlib.util.module_from_spec(_spec)
sys.modules["verify_remote_facts"] = vrf
_spec.loader.exec_module(vrf)


class VerifyRemoteFactsEncodingTests(unittest.TestCase):
    def test_run_json_decodes_utf8_chinese(self):
        """run_json must decode UTF-8 Chinese JSON output (Windows GBK fix #479)."""
        chinese_json = json.dumps({"title": "修复中文测试", "state": "closed"})
        mock_result = mock.Mock(returncode=0, stdout=chinese_json, stderr="")
        with mock.patch("verify_remote_facts.subprocess.run", return_value=mock_result) as run:
            data = vrf.run_json(["./gc", "issue", "view", "1", "--json"])
        self.assertEqual(data["title"], "修复中文测试")
        _, kwargs = run.call_args
        self.assertEqual(kwargs.get("encoding"), "utf-8")

    def test_run_json_raises_on_nonzero_exit(self):
        mock_result = mock.Mock(returncode=1, stdout="", stderr="boom")
        with mock.patch("verify_remote_facts.subprocess.run", return_value=mock_result):
            with self.assertRaises(SystemExit):
                vrf.run_json(["./gc", "issue", "view", "1", "--json"])

    def test_merged_in_main_uses_utf8_encoding(self):
        """merged_in_main must pass encoding=utf-8 (#479)."""
        mock_result = mock.Mock(returncode=0, stdout="", stderr="")
        with mock.patch("verify_remote_facts.subprocess.run", return_value=mock_result) as run:
            self.assertTrue(vrf.merged_in_main("abc123"))
        _, kwargs = run.call_args
        self.assertEqual(kwargs.get("encoding"), "utf-8")


if __name__ == "__main__":
    unittest.main()
