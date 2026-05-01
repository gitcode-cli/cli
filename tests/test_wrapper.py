import os
import runpy
import stat
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from gc_cli import wrapper


class WrapperTests(unittest.TestCase):
    def test_windows_binary_name(self):
        with mock.patch("platform.system", return_value="Windows"):
            with mock.patch("platform.machine", return_value="AMD64"):
                self.assertEqual(wrapper.get_binary_name(), "gc-windows-amd64.exe")

    def test_ensure_executable_skips_windows_execute_bit_check(self):
        with tempfile.NamedTemporaryFile() as tmp:
            path = Path(tmp.name)
            with mock.patch("platform.system", return_value="Windows"):
                with mock.patch("os.access", side_effect=AssertionError("should not check X_OK")):
                    wrapper.ensure_executable(path)

    def test_ensure_executable_chmods_non_executable_posix_binary(self):
        with tempfile.NamedTemporaryFile() as tmp:
            path = Path(tmp.name)
            path.chmod(stat.S_IRUSR | stat.S_IWUSR)

            with mock.patch("platform.system", return_value="Linux"):
                wrapper.ensure_executable(path)

            self.assertTrue(os.access(path, os.X_OK))

    def test_main_runs_packaged_binary_with_arguments(self):
        with tempfile.NamedTemporaryFile() as tmp:
            binary = Path(tmp.name)
            completed = mock.Mock(returncode=7)

            with mock.patch("gc_cli.wrapper.get_binary_path", return_value=binary):
                with mock.patch("gc_cli.wrapper.ensure_executable") as ensure_executable:
                    with mock.patch("subprocess.run", return_value=completed) as run:
                        with mock.patch("sys.argv", ["gitcode", "version"]):
                            self.assertEqual(wrapper.main(), 7)

            ensure_executable.assert_called_once_with(binary)
            run.assert_called_once_with([str(binary), "version"], cwd=os.getcwd())

    def test_module_entrypoint_delegates_to_wrapper_main(self):
        with mock.patch("gc_cli.wrapper.main", return_value=0) as main:
            with self.assertRaises(SystemExit) as exc:
                runpy.run_module("gc_cli", run_name="__main__")

        self.assertEqual(exc.exception.code, 0)
        main.assert_called_once_with()


if __name__ == "__main__":
    unittest.main()
