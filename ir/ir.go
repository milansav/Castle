package ir

type InstructionType uint

const (
	IT_NOOP      InstructionType = iota // Do nothing
	IT_SCOPE                            // Start of scope
	IT_END                              // End of scope
	IT_JMP                              // Go to scope
	IT_BE_CALL                          // Call backend api - print, etc
	IT_DEF_STACK                        // Define stack
)
