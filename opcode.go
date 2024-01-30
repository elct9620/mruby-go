package mruby

type opCode = uint8

//nolint:unused
const (
	opNop opCode = iota
	opMove
	opLoadL
	opLoadI
	opLoadINeg
	opLoadI__1
	opLoadI_0
	opLoadI_1
	opLoadI_2
	opLoadI_3
	opLoadI_4
	opLoadI_5
	opLoadI_6
	opLoadI_7
	opLoadI16
	opLoadI32
	opLoadSym
	opLoadNil
	opLoadSelf
	opLoadT
	opLoadF
	opGetGV
	opSetGV
	opGetSV
	opSetSV
	opGetIV
	opSetIV
	opGetCV
	opSetCV
	opGetConst
	opSetConst
	opGetMCnst
	opSetMCnst
	opGetUpVar
	opSetUpVar
	opGetIdx
	opSetIdx
	opJmp
	opJmpIf
	opJmpNot
	opJmpNil
	opJmpUW
	opExpect
	opRescue
	opRaiseIf
	opSelfSend
	opSelfSendB
	opSend
	opSendB
	opCall
	opSuper
	opARGARY
	opEnter
	opKey_P
	opKeyEnd
	opKArg
	opReturn
	opReturn_Blk
	opBreak
	opBlkPush
	opAdd
	opAddI
	opSub
	opSubI
	opMul
	opDiv
	opEQ
	opLT
	opLE
	opGT
	opGE
	opArray
	opArray2
	opAryCat
	opAryPush
	opArySplat
	opARef
	opASet
	opAPost
	opIntern
	opSymbol
	opString
	opStrCat
	opHash
	opHashAdd
	opHashCat
	opLambda
	opBlock
	opMethod
	opRange_Inc
	opRange_Exc
	opOClass
	opClass
	opModule
	opExec
	opDef
	opAlias
	opUndef
	opSClass
	opTClass
	opDebug
	opErr
	opExt1
	opExt2
	opExt3
	opStop
)
