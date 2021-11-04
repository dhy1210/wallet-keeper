package omnijson

type OmniCreateRawTxReferenceResult = string

type OmniCreateRawTxReferenceCommand struct {
	Raw         string
	Destination string
	//Amount      string
}

func (OmniCreateRawTxReferenceCommand) Method() string {
	return "omni_createrawtx_reference"
}

func (OmniCreateRawTxReferenceCommand) ID() string {
	return "1"
}

func (cmd OmniCreateRawTxReferenceCommand) Params() []interface{} {
	//return []interface{}{cmd.Raw, cmd.Destination, cmd.Amount}
	return []interface{}{cmd.Raw, cmd.Destination}
}
