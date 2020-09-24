# REQUIRED SECTION
include ./.defaults.mk
# END OF REQUIRED SECTION

LD_FLAGS := -s -w -X $(MODULE)/cmd.version=$(LD_VERSION) -X $(MODULE)/cmd.commit=$(LD_COMMIT) -X $(MODULE)/cmd.date=$(LD_DATE)