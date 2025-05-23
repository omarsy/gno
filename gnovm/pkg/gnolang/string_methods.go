// Code generated by "stringer -type=Kind,Op,TransCtrl,TransField,VPType,Word -output string_methods.go ."; DO NOT EDIT.

package gnolang

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidKind-0]
	_ = x[BoolKind-1]
	_ = x[StringKind-2]
	_ = x[IntKind-3]
	_ = x[Int8Kind-4]
	_ = x[Int16Kind-5]
	_ = x[Int32Kind-6]
	_ = x[Int64Kind-7]
	_ = x[UintKind-8]
	_ = x[Uint8Kind-9]
	_ = x[Uint16Kind-10]
	_ = x[Uint32Kind-11]
	_ = x[Uint64Kind-12]
	_ = x[Float32Kind-13]
	_ = x[Float64Kind-14]
	_ = x[BigintKind-15]
	_ = x[BigdecKind-16]
	_ = x[ArrayKind-17]
	_ = x[SliceKind-18]
	_ = x[PointerKind-19]
	_ = x[StructKind-20]
	_ = x[PackageKind-21]
	_ = x[InterfaceKind-22]
	_ = x[ChanKind-23]
	_ = x[FuncKind-24]
	_ = x[MapKind-25]
	_ = x[TypeKind-26]
	_ = x[BlockKind-27]
	_ = x[HeapItemKind-28]
	_ = x[TupleKind-29]
	_ = x[RefTypeKind-30]
}

const _Kind_name = "InvalidKindBoolKindStringKindIntKindInt8KindInt16KindInt32KindInt64KindUintKindUint8KindUint16KindUint32KindUint64KindFloat32KindFloat64KindBigintKindBigdecKindArrayKindSliceKindPointerKindStructKindPackageKindInterfaceKindChanKindFuncKindMapKindTypeKindBlockKindHeapItemKindTupleKindRefTypeKind"

var _Kind_index = [...]uint16{0, 11, 19, 29, 36, 44, 53, 62, 71, 79, 88, 98, 108, 118, 129, 140, 150, 160, 169, 178, 189, 199, 210, 223, 231, 239, 246, 254, 263, 275, 284, 295}

func (i Kind) String() string {
	if i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpInvalid-0]
	_ = x[OpHalt-1]
	_ = x[OpNoop-2]
	_ = x[OpExec-3]
	_ = x[OpPrecall-4]
	_ = x[OpCall-5]
	_ = x[OpCallNativeBody-6]
	_ = x[OpDefer-10]
	_ = x[OpCallDeferNativeBody-11]
	_ = x[OpGo-12]
	_ = x[OpSelect-13]
	_ = x[OpSwitchClause-14]
	_ = x[OpSwitchClauseCase-15]
	_ = x[OpTypeSwitch-16]
	_ = x[OpIfCond-17]
	_ = x[OpPopValue-18]
	_ = x[OpPopResults-19]
	_ = x[OpPopBlock-20]
	_ = x[OpPopFrameAndReset-21]
	_ = x[OpPanic1-22]
	_ = x[OpPanic2-23]
	_ = x[OpReturn-26]
	_ = x[OpReturnAfterCopy-27]
	_ = x[OpReturnFromBlock-28]
	_ = x[OpReturnToBlock-29]
	_ = x[OpUpos-32]
	_ = x[OpUneg-33]
	_ = x[OpUnot-34]
	_ = x[OpUxor-35]
	_ = x[OpUrecv-37]
	_ = x[OpLor-38]
	_ = x[OpLand-39]
	_ = x[OpEql-40]
	_ = x[OpNeq-41]
	_ = x[OpLss-42]
	_ = x[OpLeq-43]
	_ = x[OpGtr-44]
	_ = x[OpGeq-45]
	_ = x[OpAdd-46]
	_ = x[OpSub-47]
	_ = x[OpBor-48]
	_ = x[OpXor-49]
	_ = x[OpMul-50]
	_ = x[OpQuo-51]
	_ = x[OpRem-52]
	_ = x[OpShl-53]
	_ = x[OpShr-54]
	_ = x[OpBand-55]
	_ = x[OpBandn-56]
	_ = x[OpEval-64]
	_ = x[OpBinary1-65]
	_ = x[OpIndex1-66]
	_ = x[OpIndex2-67]
	_ = x[OpSelector-68]
	_ = x[OpSlice-69]
	_ = x[OpStar-70]
	_ = x[OpRef-71]
	_ = x[OpTypeAssert1-72]
	_ = x[OpTypeAssert2-73]
	_ = x[OpStaticTypeOf-74]
	_ = x[OpCompositeLit-75]
	_ = x[OpArrayLit-76]
	_ = x[OpSliceLit-77]
	_ = x[OpSliceLit2-78]
	_ = x[OpMapLit-79]
	_ = x[OpStructLit-80]
	_ = x[OpFuncLit-81]
	_ = x[OpConvert-82]
	_ = x[OpFieldType-112]
	_ = x[OpArrayType-113]
	_ = x[OpSliceType-114]
	_ = x[OpPointerType-115]
	_ = x[OpInterfaceType-116]
	_ = x[OpChanType-117]
	_ = x[OpFuncType-118]
	_ = x[OpMapType-119]
	_ = x[OpStructType-120]
	_ = x[OpAssign-128]
	_ = x[OpAddAssign-129]
	_ = x[OpSubAssign-130]
	_ = x[OpMulAssign-131]
	_ = x[OpQuoAssign-132]
	_ = x[OpRemAssign-133]
	_ = x[OpBandAssign-134]
	_ = x[OpBandnAssign-135]
	_ = x[OpBorAssign-136]
	_ = x[OpXorAssign-137]
	_ = x[OpShlAssign-138]
	_ = x[OpShrAssign-139]
	_ = x[OpDefine-140]
	_ = x[OpInc-141]
	_ = x[OpDec-142]
	_ = x[OpValueDecl-144]
	_ = x[OpTypeDecl-145]
	_ = x[OpSticky-208]
	_ = x[OpBody-209]
	_ = x[OpForLoop-210]
	_ = x[OpRangeIter-211]
	_ = x[OpRangeIterString-212]
	_ = x[OpRangeIterMap-213]
	_ = x[OpRangeIterArrayPtr-214]
	_ = x[OpReturnCallDefers-215]
	_ = x[OpVoid-255]
}

const _Op_name = "OpInvalidOpHaltOpNoopOpExecOpPrecallOpCallOpCallNativeBodyOpDeferOpCallDeferNativeBodyOpGoOpSelectOpSwitchClauseOpSwitchClauseCaseOpTypeSwitchOpIfCondOpPopValueOpPopResultsOpPopBlockOpPopFrameAndResetOpPanic1OpPanic2OpReturnOpReturnAfterCopyOpReturnFromBlockOpReturnToBlockOpUposOpUnegOpUnotOpUxorOpUrecvOpLorOpLandOpEqlOpNeqOpLssOpLeqOpGtrOpGeqOpAddOpSubOpBorOpXorOpMulOpQuoOpRemOpShlOpShrOpBandOpBandnOpEvalOpBinary1OpIndex1OpIndex2OpSelectorOpSliceOpStarOpRefOpTypeAssert1OpTypeAssert2OpStaticTypeOfOpCompositeLitOpArrayLitOpSliceLitOpSliceLit2OpMapLitOpStructLitOpFuncLitOpConvertOpFieldTypeOpArrayTypeOpSliceTypeOpPointerTypeOpInterfaceTypeOpChanTypeOpFuncTypeOpMapTypeOpStructTypeOpAssignOpAddAssignOpSubAssignOpMulAssignOpQuoAssignOpRemAssignOpBandAssignOpBandnAssignOpBorAssignOpXorAssignOpShlAssignOpShrAssignOpDefineOpIncOpDecOpValueDeclOpTypeDeclOpStickyOpBodyOpForLoopOpRangeIterOpRangeIterStringOpRangeIterMapOpRangeIterArrayPtrOpReturnCallDefersOpVoid"

var _Op_map = map[Op]string{
	0:   _Op_name[0:9],
	1:   _Op_name[9:15],
	2:   _Op_name[15:21],
	3:   _Op_name[21:27],
	4:   _Op_name[27:36],
	5:   _Op_name[36:42],
	6:   _Op_name[42:58],
	10:  _Op_name[58:65],
	11:  _Op_name[65:86],
	12:  _Op_name[86:90],
	13:  _Op_name[90:98],
	14:  _Op_name[98:112],
	15:  _Op_name[112:130],
	16:  _Op_name[130:142],
	17:  _Op_name[142:150],
	18:  _Op_name[150:160],
	19:  _Op_name[160:172],
	20:  _Op_name[172:182],
	21:  _Op_name[182:200],
	22:  _Op_name[200:208],
	23:  _Op_name[208:216],
	26:  _Op_name[216:224],
	27:  _Op_name[224:241],
	28:  _Op_name[241:258],
	29:  _Op_name[258:273],
	32:  _Op_name[273:279],
	33:  _Op_name[279:285],
	34:  _Op_name[285:291],
	35:  _Op_name[291:297],
	37:  _Op_name[297:304],
	38:  _Op_name[304:309],
	39:  _Op_name[309:315],
	40:  _Op_name[315:320],
	41:  _Op_name[320:325],
	42:  _Op_name[325:330],
	43:  _Op_name[330:335],
	44:  _Op_name[335:340],
	45:  _Op_name[340:345],
	46:  _Op_name[345:350],
	47:  _Op_name[350:355],
	48:  _Op_name[355:360],
	49:  _Op_name[360:365],
	50:  _Op_name[365:370],
	51:  _Op_name[370:375],
	52:  _Op_name[375:380],
	53:  _Op_name[380:385],
	54:  _Op_name[385:390],
	55:  _Op_name[390:396],
	56:  _Op_name[396:403],
	64:  _Op_name[403:409],
	65:  _Op_name[409:418],
	66:  _Op_name[418:426],
	67:  _Op_name[426:434],
	68:  _Op_name[434:444],
	69:  _Op_name[444:451],
	70:  _Op_name[451:457],
	71:  _Op_name[457:462],
	72:  _Op_name[462:475],
	73:  _Op_name[475:488],
	74:  _Op_name[488:502],
	75:  _Op_name[502:516],
	76:  _Op_name[516:526],
	77:  _Op_name[526:536],
	78:  _Op_name[536:547],
	79:  _Op_name[547:555],
	80:  _Op_name[555:566],
	81:  _Op_name[566:575],
	82:  _Op_name[575:584],
	112: _Op_name[584:595],
	113: _Op_name[595:606],
	114: _Op_name[606:617],
	115: _Op_name[617:630],
	116: _Op_name[630:645],
	117: _Op_name[645:655],
	118: _Op_name[655:665],
	119: _Op_name[665:674],
	120: _Op_name[674:686],
	128: _Op_name[686:694],
	129: _Op_name[694:705],
	130: _Op_name[705:716],
	131: _Op_name[716:727],
	132: _Op_name[727:738],
	133: _Op_name[738:749],
	134: _Op_name[749:761],
	135: _Op_name[761:774],
	136: _Op_name[774:785],
	137: _Op_name[785:796],
	138: _Op_name[796:807],
	139: _Op_name[807:818],
	140: _Op_name[818:826],
	141: _Op_name[826:831],
	142: _Op_name[831:836],
	144: _Op_name[836:847],
	145: _Op_name[847:857],
	208: _Op_name[857:865],
	209: _Op_name[865:871],
	210: _Op_name[871:880],
	211: _Op_name[880:891],
	212: _Op_name[891:908],
	213: _Op_name[908:922],
	214: _Op_name[922:941],
	215: _Op_name[941:959],
	255: _Op_name[959:965],
}

func (i Op) String() string {
	if str, ok := _Op_map[i]; ok {
		return str
	}
	return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TRANS_CONTINUE-0]
	_ = x[TRANS_SKIP-1]
	_ = x[TRANS_EXIT-2]
}

const _TransCtrl_name = "TRANS_CONTINUETRANS_SKIPTRANS_EXIT"

var _TransCtrl_index = [...]uint8{0, 14, 24, 34}

func (i TransCtrl) String() string {
	if i >= TransCtrl(len(_TransCtrl_index)-1) {
		return "TransCtrl(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransCtrl_name[_TransCtrl_index[i]:_TransCtrl_index[i+1]]
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TRANS_ROOT-0]
	_ = x[TRANS_BINARY_LEFT-1]
	_ = x[TRANS_BINARY_RIGHT-2]
	_ = x[TRANS_CALL_FUNC-3]
	_ = x[TRANS_CALL_ARG-4]
	_ = x[TRANS_INDEX_X-5]
	_ = x[TRANS_INDEX_INDEX-6]
	_ = x[TRANS_SELECTOR_X-7]
	_ = x[TRANS_SLICE_X-8]
	_ = x[TRANS_SLICE_LOW-9]
	_ = x[TRANS_SLICE_HIGH-10]
	_ = x[TRANS_SLICE_MAX-11]
	_ = x[TRANS_STAR_X-12]
	_ = x[TRANS_REF_X-13]
	_ = x[TRANS_TYPEASSERT_X-14]
	_ = x[TRANS_TYPEASSERT_TYPE-15]
	_ = x[TRANS_UNARY_X-16]
	_ = x[TRANS_COMPOSITE_TYPE-17]
	_ = x[TRANS_COMPOSITE_KEY-18]
	_ = x[TRANS_COMPOSITE_VALUE-19]
	_ = x[TRANS_FUNCLIT_TYPE-20]
	_ = x[TRANS_FUNCLIT_HEAP_CAPTURE-21]
	_ = x[TRANS_FUNCLIT_BODY-22]
	_ = x[TRANS_FIELDTYPE_NAME-23]
	_ = x[TRANS_FIELDTYPE_TYPE-24]
	_ = x[TRANS_FIELDTYPE_TAG-25]
	_ = x[TRANS_ARRAYTYPE_LEN-26]
	_ = x[TRANS_ARRAYTYPE_ELT-27]
	_ = x[TRANS_SLICETYPE_ELT-28]
	_ = x[TRANS_INTERFACETYPE_METHOD-29]
	_ = x[TRANS_CHANTYPE_VALUE-30]
	_ = x[TRANS_FUNCTYPE_PARAM-31]
	_ = x[TRANS_FUNCTYPE_RESULT-32]
	_ = x[TRANS_MAPTYPE_KEY-33]
	_ = x[TRANS_MAPTYPE_VALUE-34]
	_ = x[TRANS_STRUCTTYPE_FIELD-35]
	_ = x[TRANS_ASSIGN_LHS-36]
	_ = x[TRANS_ASSIGN_RHS-37]
	_ = x[TRANS_BLOCK_BODY-38]
	_ = x[TRANS_DECL_BODY-39]
	_ = x[TRANS_DEFER_CALL-40]
	_ = x[TRANS_EXPR_X-41]
	_ = x[TRANS_FOR_INIT-42]
	_ = x[TRANS_FOR_COND-43]
	_ = x[TRANS_FOR_POST-44]
	_ = x[TRANS_FOR_BODY-45]
	_ = x[TRANS_GO_CALL-46]
	_ = x[TRANS_IF_INIT-47]
	_ = x[TRANS_IF_COND-48]
	_ = x[TRANS_IF_BODY-49]
	_ = x[TRANS_IF_ELSE-50]
	_ = x[TRANS_IF_CASE_BODY-51]
	_ = x[TRANS_INCDEC_X-52]
	_ = x[TRANS_RANGE_X-53]
	_ = x[TRANS_RANGE_KEY-54]
	_ = x[TRANS_RANGE_VALUE-55]
	_ = x[TRANS_RANGE_BODY-56]
	_ = x[TRANS_RETURN_RESULT-57]
	_ = x[TRANS_SELECT_CASE-58]
	_ = x[TRANS_SELECTCASE_COMM-59]
	_ = x[TRANS_SELECTCASE_BODY-60]
	_ = x[TRANS_SEND_CHAN-61]
	_ = x[TRANS_SEND_VALUE-62]
	_ = x[TRANS_SWITCH_INIT-63]
	_ = x[TRANS_SWITCH_X-64]
	_ = x[TRANS_SWITCH_CASE-65]
	_ = x[TRANS_SWITCHCASE_CASE-66]
	_ = x[TRANS_SWITCHCASE_BODY-67]
	_ = x[TRANS_FUNC_RECV-68]
	_ = x[TRANS_FUNC_TYPE-69]
	_ = x[TRANS_FUNC_BODY-70]
	_ = x[TRANS_IMPORT_PATH-71]
	_ = x[TRANS_CONST_TYPE-72]
	_ = x[TRANS_CONST_VALUE-73]
	_ = x[TRANS_VAR_NAME-74]
	_ = x[TRANS_VAR_TYPE-75]
	_ = x[TRANS_VAR_VALUE-76]
	_ = x[TRANS_TYPE_TYPE-77]
	_ = x[TRANS_FILE_BODY-78]
}

const _TransField_name = "TRANS_ROOTTRANS_BINARY_LEFTTRANS_BINARY_RIGHTTRANS_CALL_FUNCTRANS_CALL_ARGTRANS_INDEX_XTRANS_INDEX_INDEXTRANS_SELECTOR_XTRANS_SLICE_XTRANS_SLICE_LOWTRANS_SLICE_HIGHTRANS_SLICE_MAXTRANS_STAR_XTRANS_REF_XTRANS_TYPEASSERT_XTRANS_TYPEASSERT_TYPETRANS_UNARY_XTRANS_COMPOSITE_TYPETRANS_COMPOSITE_KEYTRANS_COMPOSITE_VALUETRANS_FUNCLIT_TYPETRANS_FUNCLIT_HEAP_CAPTURETRANS_FUNCLIT_BODYTRANS_FIELDTYPE_NAMETRANS_FIELDTYPE_TYPETRANS_FIELDTYPE_TAGTRANS_ARRAYTYPE_LENTRANS_ARRAYTYPE_ELTTRANS_SLICETYPE_ELTTRANS_INTERFACETYPE_METHODTRANS_CHANTYPE_VALUETRANS_FUNCTYPE_PARAMTRANS_FUNCTYPE_RESULTTRANS_MAPTYPE_KEYTRANS_MAPTYPE_VALUETRANS_STRUCTTYPE_FIELDTRANS_ASSIGN_LHSTRANS_ASSIGN_RHSTRANS_BLOCK_BODYTRANS_DECL_BODYTRANS_DEFER_CALLTRANS_EXPR_XTRANS_FOR_INITTRANS_FOR_CONDTRANS_FOR_POSTTRANS_FOR_BODYTRANS_GO_CALLTRANS_IF_INITTRANS_IF_CONDTRANS_IF_BODYTRANS_IF_ELSETRANS_IF_CASE_BODYTRANS_INCDEC_XTRANS_RANGE_XTRANS_RANGE_KEYTRANS_RANGE_VALUETRANS_RANGE_BODYTRANS_RETURN_RESULTTRANS_SELECT_CASETRANS_SELECTCASE_COMMTRANS_SELECTCASE_BODYTRANS_SEND_CHANTRANS_SEND_VALUETRANS_SWITCH_INITTRANS_SWITCH_XTRANS_SWITCH_CASETRANS_SWITCHCASE_CASETRANS_SWITCHCASE_BODYTRANS_FUNC_RECVTRANS_FUNC_TYPETRANS_FUNC_BODYTRANS_IMPORT_PATHTRANS_CONST_TYPETRANS_CONST_VALUETRANS_VAR_NAMETRANS_VAR_TYPETRANS_VAR_VALUETRANS_TYPE_TYPETRANS_FILE_BODY"

var _TransField_index = [...]uint16{0, 10, 27, 45, 60, 74, 87, 104, 120, 133, 148, 164, 179, 191, 202, 220, 241, 254, 274, 293, 314, 332, 358, 376, 396, 416, 435, 454, 473, 492, 518, 538, 558, 579, 596, 615, 637, 653, 669, 685, 700, 716, 728, 742, 756, 770, 784, 797, 810, 823, 836, 849, 867, 881, 894, 909, 926, 942, 961, 978, 999, 1020, 1035, 1051, 1068, 1082, 1099, 1120, 1141, 1156, 1171, 1186, 1203, 1219, 1236, 1250, 1264, 1279, 1294, 1309}

func (i TransField) String() string {
	if i >= TransField(len(_TransField_index)-1) {
		return "TransField(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransField_name[_TransField_index[i]:_TransField_index[i+1]]
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[VPUverse-0]
	_ = x[VPBlock-1]
	_ = x[VPField-2]
	_ = x[VPValMethod-3]
	_ = x[VPPtrMethod-4]
	_ = x[VPInterface-5]
	_ = x[VPSubrefField-6]
	_ = x[VPDerefField-18]
	_ = x[VPDerefValMethod-19]
	_ = x[VPDerefPtrMethod-20]
	_ = x[VPDerefInterface-21]
}

const (
	_VPType_name_0 = "VPUverseVPBlockVPFieldVPValMethodVPPtrMethodVPInterfaceVPSubrefField"
	_VPType_name_1 = "VPDerefFieldVPDerefValMethodVPDerefPtrMethodVPDerefInterface"
)

var (
	_VPType_index_0 = [...]uint8{0, 8, 15, 22, 33, 44, 55, 68}
	_VPType_index_1 = [...]uint8{0, 12, 28, 44, 60}
)

func (i VPType) String() string {
	switch {
	case i <= 6:
		return _VPType_name_0[_VPType_index_0[i]:_VPType_index_0[i+1]]
	case 18 <= i && i <= 21:
		i -= 18
		return _VPType_name_1[_VPType_index_1[i]:_VPType_index_1[i+1]]
	default:
		return "VPType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[NAME-1]
	_ = x[INT-2]
	_ = x[FLOAT-3]
	_ = x[IMAG-4]
	_ = x[CHAR-5]
	_ = x[STRING-6]
	_ = x[ADD-7]
	_ = x[SUB-8]
	_ = x[MUL-9]
	_ = x[QUO-10]
	_ = x[REM-11]
	_ = x[BAND-12]
	_ = x[BOR-13]
	_ = x[XOR-14]
	_ = x[SHL-15]
	_ = x[SHR-16]
	_ = x[BAND_NOT-17]
	_ = x[ADD_ASSIGN-18]
	_ = x[SUB_ASSIGN-19]
	_ = x[MUL_ASSIGN-20]
	_ = x[QUO_ASSIGN-21]
	_ = x[REM_ASSIGN-22]
	_ = x[BAND_ASSIGN-23]
	_ = x[BOR_ASSIGN-24]
	_ = x[XOR_ASSIGN-25]
	_ = x[SHL_ASSIGN-26]
	_ = x[SHR_ASSIGN-27]
	_ = x[BAND_NOT_ASSIGN-28]
	_ = x[LAND-29]
	_ = x[LOR-30]
	_ = x[ARROW-31]
	_ = x[INC-32]
	_ = x[DEC-33]
	_ = x[EQL-34]
	_ = x[LSS-35]
	_ = x[GTR-36]
	_ = x[ASSIGN-37]
	_ = x[NOT-38]
	_ = x[NEQ-39]
	_ = x[LEQ-40]
	_ = x[GEQ-41]
	_ = x[DEFINE-42]
	_ = x[BREAK-43]
	_ = x[CASE-44]
	_ = x[CHAN-45]
	_ = x[CONST-46]
	_ = x[CONTINUE-47]
	_ = x[DEFAULT-48]
	_ = x[DEFER-49]
	_ = x[ELSE-50]
	_ = x[FALLTHROUGH-51]
	_ = x[FOR-52]
	_ = x[FUNC-53]
	_ = x[GO-54]
	_ = x[GOTO-55]
	_ = x[IF-56]
	_ = x[IMPORT-57]
	_ = x[INTERFACE-58]
	_ = x[MAP-59]
	_ = x[PACKAGE-60]
	_ = x[RANGE-61]
	_ = x[RETURN-62]
	_ = x[SELECT-63]
	_ = x[STRUCT-64]
	_ = x[SWITCH-65]
	_ = x[TYPE-66]
	_ = x[VAR-67]
}

const _Word_name = "ILLEGALNAMEINTFLOATIMAGCHARSTRINGADDSUBMULQUOREMBANDBORXORSHLSHRBAND_NOTADD_ASSIGNSUB_ASSIGNMUL_ASSIGNQUO_ASSIGNREM_ASSIGNBAND_ASSIGNBOR_ASSIGNXOR_ASSIGNSHL_ASSIGNSHR_ASSIGNBAND_NOT_ASSIGNLANDLORARROWINCDECEQLLSSGTRASSIGNNOTNEQLEQGEQDEFINEBREAKCASECHANCONSTCONTINUEDEFAULTDEFERELSEFALLTHROUGHFORFUNCGOGOTOIFIMPORTINTERFACEMAPPACKAGERANGERETURNSELECTSTRUCTSWITCHTYPEVAR"

var _Word_index = [...]uint16{0, 7, 11, 14, 19, 23, 27, 33, 36, 39, 42, 45, 48, 52, 55, 58, 61, 64, 72, 82, 92, 102, 112, 122, 133, 143, 153, 163, 173, 188, 192, 195, 200, 203, 206, 209, 212, 215, 221, 224, 227, 230, 233, 239, 244, 248, 252, 257, 265, 272, 277, 281, 292, 295, 299, 301, 305, 307, 313, 322, 325, 332, 337, 343, 349, 355, 361, 365, 368}

func (i Word) String() string {
	if i < 0 || i >= Word(len(_Word_index)-1) {
		return "Word(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Word_name[_Word_index[i]:_Word_index[i+1]]
}
