// Code generated by "stringer -trimprefix Op -type Opcode opcodes.go"; DO NOT EDIT.

package vm

import "strconv"

const _Opcode_name = "NopDropDrop2DupDup2SwapOverPickRollRetFailZeroPush1Push2Push3Push4Push5Push6Push7Push8Push64OneNeg1PushTNowRandPushLAddSubMulDivModNotNegIncDecIndexLenAppendExtendSliceFieldFieldLIfzIfnzElseEndSumAvgMaxMinChoiceWChoiceSortLookup"

var _Opcode_map = map[Opcode]string{
	0:   _Opcode_name[0:3],
	1:   _Opcode_name[3:7],
	2:   _Opcode_name[7:12],
	5:   _Opcode_name[12:15],
	6:   _Opcode_name[15:19],
	9:   _Opcode_name[19:23],
	13:  _Opcode_name[23:27],
	14:  _Opcode_name[27:31],
	15:  _Opcode_name[31:35],
	16:  _Opcode_name[35:38],
	17:  _Opcode_name[38:42],
	32:  _Opcode_name[42:46],
	33:  _Opcode_name[46:51],
	34:  _Opcode_name[51:56],
	35:  _Opcode_name[56:61],
	36:  _Opcode_name[61:66],
	37:  _Opcode_name[66:71],
	38:  _Opcode_name[71:76],
	39:  _Opcode_name[76:81],
	40:  _Opcode_name[81:86],
	41:  _Opcode_name[86:92],
	42:  _Opcode_name[92:95],
	43:  _Opcode_name[95:99],
	44:  _Opcode_name[99:104],
	45:  _Opcode_name[104:107],
	47:  _Opcode_name[107:111],
	48:  _Opcode_name[111:116],
	64:  _Opcode_name[116:119],
	65:  _Opcode_name[119:122],
	66:  _Opcode_name[122:125],
	67:  _Opcode_name[125:128],
	68:  _Opcode_name[128:131],
	69:  _Opcode_name[131:134],
	70:  _Opcode_name[134:137],
	71:  _Opcode_name[137:140],
	72:  _Opcode_name[140:143],
	80:  _Opcode_name[143:148],
	81:  _Opcode_name[148:151],
	82:  _Opcode_name[151:157],
	83:  _Opcode_name[157:163],
	84:  _Opcode_name[163:168],
	96:  _Opcode_name[168:173],
	112: _Opcode_name[173:179],
	128: _Opcode_name[179:182],
	129: _Opcode_name[182:186],
	130: _Opcode_name[186:190],
	136: _Opcode_name[190:193],
	144: _Opcode_name[193:196],
	145: _Opcode_name[196:199],
	146: _Opcode_name[199:202],
	147: _Opcode_name[202:205],
	148: _Opcode_name[205:211],
	149: _Opcode_name[211:218],
	150: _Opcode_name[218:222],
	151: _Opcode_name[222:228],
}

func (i Opcode) String() string {
	if str, ok := _Opcode_map[i]; ok {
		return str
	}
	return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
}
