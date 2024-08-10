import os
import subprocess
import filecmp
import pytest
import shutil
import time

@pytest.fixture(scope="function")
def render_setup():
    if os.path.exists("result"):
        shutil.rmtree("result")
    os.mkdir("result")
    yield
    if os.path.exists("result"):
        shutil.rmtree("result")

def test_render_default_values(render_setup):
    command = ["./dcpm", "render", "-o", "result", "-p", "mock_project"]
    result = subprocess.run(command, capture_output=True, text=True)

    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    comparison = filecmp.dircmp("result", "expected1")
    if comparison.left_only or comparison.right_only or comparison.diff_files:
        for name in comparison.left_only:
            print(f"Only in result: {name}")
        for name in comparison.right_only:
            print(f"Only in expected result: {name}")
        for name in comparison.diff_files:
            print(f"Different files: {name}")
        assert False, "Directories do not match"
    

def test_render_custom_values(render_setup):
    command = ["./dcpm", "render", "-o", "result", "-p", "mock_project", "-v", "values2.yaml"]
    result = subprocess.run(command, capture_output=True, text=True)

    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    comparison = filecmp.dircmp("result", "expected2")
    if comparison.left_only or comparison.right_only or comparison.diff_files:
        for name in comparison.left_only:
            print(f"Only in result: {name}")
        for name in comparison.right_only:
            print(f"Only in expected result: {name}")
        for name in comparison.diff_files:
            print(f"Different files: {name}")
        assert False, "Directories do not match"


@pytest.fixture(scope="function")
def render_errors_setup():
    if os.path.exists("result"):
        shutil.rmtree("result")
    os.mkdir("result")
    if os.path.exists("empty_project_test"):
        shutil.rmtree("empty_project_test")
    os.mkdir("empty_project_test")
    yield
    if os.path.exists("result"):
        shutil.rmtree("result")
    if os.path.exists("empty_project_test"):
        shutil.rmtree("empty_project_test")

def test_render_errors(render_errors_setup):
    command = ["./dcpm", "render", "-o", "result", "-p", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)

    assert result.returncode == 1
    assert "can't find values.yaml inside direcotry:" in result.stderr
    assert "can't find templates direcotry inside direcotry:" in result.stderr
    assert "can't find running_config direcotry inside direcotry:" in result.stderr
    assert "can't find dependencies direcotry inside direcotry:" in result.stderr


@pytest.fixture(scope="function")
def checksum_setup():
    if os.path.exists("empty_project_test"):
        shutil.rmtree("empty_project_test")
    shutil.copytree("new_empty_project", "empty_project_test")
    yield
    if os.path.exists("empty_project_test"):
        shutil.rmtree("empty_project_test")

def test_checksum_after_new_file(checksum_setup):
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"
   
   # Create new file in the project directory
    with open('empty_project_test/new_file.txt', 'w') as file:
        pass
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"

    # Remove a file
    os.remove("empty_project_test/new_file.txt")
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"



def test_checksum_after_new_directory(checksum_setup):
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"

    # Create new directory
    os.makedirs("./empty_project_test/new_dir")
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"

    # Remove directory 
    shutil.rmtree("./empty_project_test/new_dir")
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"


def test_checksum_after_edited_file(checksum_setup):
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"
   
   # Edit a file
    with open('empty_project_test/values.yaml', 'w') as file:
        file.write("key: value\n")
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"

    # Edit a file
    with open('empty_project_test/values.yaml', 'w') as file:
        file.write("")
    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert result.stderr == "checksum doesn't match!"

    # Recalculate checksum
    command = ["./dcpm", "checksum", "create", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    command = ["./dcpm", "checksum", "check", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"
    assert result.stdout == "checksum is correct"

def test_checksum_errors(checksum_setup):
    # No path provided
    command = ["./dcpm", "checksum", "create"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1

    command = ["./dcpm", "checksum", "check"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1

    # Wrong path provided
    command = ["./dcpm", "checksum", "create", "non_existing_dir"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1

    command = ["./dcpm", "checksum", "check", "non_existing_dir"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1

    # Wrong path provided
    command = ["./dcpm", "checksum", "create", "non_existing_dir"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert "can't find a direcory:" in result.stderr

    command = ["./dcpm", "checksum", "check", "non_existing_dir"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert "can't find a direcory:" in result.stderr

    # No CHECKSUM file
    command = ["./dcpm", "checksum", "check", "mock_project"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    assert "can't find a CHECKSUM file:" in result.stderr

@pytest.fixture(scope="function")
def init_setup():
    yield
    if os.path.exists("empty_project_test"):
        shutil.rmtree("empty_project_test")

def test_init(init_setup):
    command = ["./dcpm", "init", "empty_project_test"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    comparison = filecmp.dircmp("empty_project_test", "new_empty_project")
    if comparison.left_only or comparison.right_only or comparison.diff_files:
        for name in comparison.left_only:
            print(f"Only in result: {name}")
        for name in comparison.right_only:
            print(f"Only in expected result: {name}")
        for name in comparison.diff_files:
            print(f"Different files: {name}")
        assert False, "Directories do not match"

def test_init_error(init_setup):
    command = ["./dcpm", "init"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 1
    

@pytest.fixture(scope="function")
def install_setup():
    if os.path.exists("new_mock_project"):
        shutil.rmtree("new_mock_project")
    shutil.copytree("mock_project", "new_mock_project")
    yield
    if os.path.exists("new_mock_project"):
        shutil.rmtree("new_mock_project")

def test_install_default_values(install_setup):
    expected_containers = ["web", "db", "redis", "phpmyadmin"]

    command = ["./dcpm", "install", "new_mock_project"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    docker_ps_command = ["docker", "ps", "--format", "{{.Names}}"]
    result = subprocess.run(docker_ps_command, capture_output=True, text=True)
    assert result.returncode == 0, f"docker ps failed with output: {result.stderr}"

    running_containers = result.stdout.splitlines()
    assert len(expected_containers) == len(running_containers), "Different number of running and expected containers"
    for i, exp_cont in enumerate(expected_containers):
        assert exp_cont in running_containers, f"{exp_cont} is not running."
        assert running_containers[i] in expected_containers, f"{running_containers[i]} in not supposed to be running"

    # Uninstall
    command = ["./dcpm", "uninstall", "new_mock_project"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    docker_ps_command = ["docker", "ps", "--format", "{{.Names}}"]
    result = subprocess.run(docker_ps_command, capture_output=True, text=True)
    assert result.returncode == 0, f"docker ps failed with output: {result.stderr}"

    running_containers = result.stdout.splitlines()
    assert len(running_containers) == 0, f"the following containers are still running after uninstall {running_containers}"


def test_install_custom_values(install_setup):
    expected_containers = ["nginx_edit", "db_edit", "phpmyadmin_edit"]

    command = ["./dcpm", "install", "new_mock_project", "-v", "values-install.yaml"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    docker_ps_command = ["docker", "ps", "--format", "{{.Names}}"]
    result = subprocess.run(docker_ps_command, capture_output=True, text=True)
    assert result.returncode == 0, f"docker ps failed with output: {result.stderr}"

    running_containers = result.stdout.splitlines()
    assert len(expected_containers) == len(running_containers), "Different number of running and expected containers"
    for i, exp_cont in enumerate(expected_containers):
        assert exp_cont in running_containers, f"{exp_cont} is not running."
        assert running_containers[i] in expected_containers, f"{running_containers[i]} in not supposed to be running"

    # Uninstall
    command = ["./dcpm", "uninstall", "new_mock_project"]
    result = subprocess.run(command, capture_output=True, text=True)
    assert result.returncode == 0, f"Command failed with output: {result.stderr}"

    docker_ps_command = ["docker", "ps", "--format", "{{.Names}}"]
    result = subprocess.run(docker_ps_command, capture_output=True, text=True)
    assert result.returncode == 0, f"docker ps failed with output: {result.stderr}"

    running_containers = result.stdout.splitlines()
    assert len(running_containers) == 0, f"the following containers are still running after uninstall {running_containers}"
