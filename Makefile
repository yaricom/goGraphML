#
# Go parameters
#
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test -count=1
GOGET = $(GOCMD) get

# The default targets to run
#
all: test

# Run unit tests
#
test:
	$(GOTEST) -v --short ./...
