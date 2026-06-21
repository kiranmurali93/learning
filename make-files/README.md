Makefiles

make files are used to automate repeated tasks. 
Used in many languages for the same.

In this make-files folder im taking an example for golang

Steps to recreate:
- create a main.go file
- create a basic function
- Create a Makefile
- Add commands in the Makefile

To run the make file:
make <command>
eg:
`make all `

NOTES:
A MakeFile has a rule for structure
```
target: deps
    shell-command
```
Use tab separator for shell commands after target.