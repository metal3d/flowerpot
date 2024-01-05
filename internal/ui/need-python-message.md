# Specific Python version is required (max 3.11)

Sorry, we didn't find suitable version of Python.

You should verify :

- if python is installed
- if the version is not > 3.11

Generally speaking, recent distributions allow you to install several versions 
of Python without conflicting with the current version.

The executable will then be named "`python3.11`".

## Python installation

To install Python with the desired version, please change "3.11" with the
version you want to try. E.g. "3.11", "3.10"

It is commonly safe to install serveral Python versions on your system.

### Ubuntun/Debian Like, try doing:

```
sudo apt update
sudo apt install python3.11 python3.11-pip python3.11-venv
```

### Fedora/Red Hat Like

```
sudo dnf install python3.11
```

## Checking

Then, check the version with :

```
python3.11 --version
```
