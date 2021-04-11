module github.com/networklore/netrasp

go 1.13

require (
	github.com/magefile/mage v1.11.0
	golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
)

replace golang.org/x/crypto => github.com/ogenstad/crypto v0.0.0-20210308070823-6d211c1ce3d7
