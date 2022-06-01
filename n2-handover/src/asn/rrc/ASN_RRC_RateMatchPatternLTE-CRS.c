/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "NR-RRC-Definitions"
 * 	found in "asn/nr-rrc-15.6.0.asn1"
 * 	`asn1c -fcompound-names -pdu=all -findirect-choice -fno-include-deps -gen-PER -no-gen-OER -no-gen-example -D rrc`
 */

#include "ASN_RRC_RateMatchPatternLTE-CRS.h"

#include "ASN_RRC_EUTRA-MBSFN-SubframeConfigList.h"
/*
 * This type is implemented using NativeEnumerated,
 * so here we adjust the DEF accordingly.
 */
/*
 * This type is implemented using NativeEnumerated,
 * so here we adjust the DEF accordingly.
 */
/*
 * This type is implemented using NativeEnumerated,
 * so here we adjust the DEF accordingly.
 */
static int
memb_ASN_RRC_carrierFreqDL_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	long value;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	value = *(const long *)sptr;
	
	if((value >= 0 && value <= 16383)) {
		/* Constraint check succeeded */
		return 0;
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static asn_per_constraints_t asn_PER_type_ASN_RRC_carrierBandwidthDL_constr_3 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 3,  3,  0,  7 }	/* (0..7) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_nrofCRS_Ports_constr_13 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 2,  2,  0,  2 }	/* (0..2) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_v_Shift_constr_17 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 3,  3,  0,  5 }	/* (0..5) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_memb_ASN_RRC_carrierFreqDL_constr_2 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 14,  14,  0,  16383 }	/* (0..16383) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_carrierBandwidthDL_value2enum_3[] = {
	{ 0,	2,	"n6" },
	{ 1,	3,	"n15" },
	{ 2,	3,	"n25" },
	{ 3,	3,	"n50" },
	{ 4,	3,	"n75" },
	{ 5,	4,	"n100" },
	{ 6,	6,	"spare2" },
	{ 7,	6,	"spare1" }
};
static const unsigned int asn_MAP_ASN_RRC_carrierBandwidthDL_enum2value_3[] = {
	5,	/* n100(5) */
	1,	/* n15(1) */
	2,	/* n25(2) */
	3,	/* n50(3) */
	0,	/* n6(0) */
	4,	/* n75(4) */
	7,	/* spare1(7) */
	6	/* spare2(6) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_carrierBandwidthDL_specs_3 = {
	asn_MAP_ASN_RRC_carrierBandwidthDL_value2enum_3,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_carrierBandwidthDL_enum2value_3,	/* N => "tag"; sorted by N */
	8,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3[] = {
	(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_carrierBandwidthDL_3 = {
	"carrierBandwidthDL",
	"carrierBandwidthDL",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3,
	sizeof(asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3)
		/sizeof(asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3)
		/sizeof(asn_DEF_ASN_RRC_carrierBandwidthDL_tags_3[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_carrierBandwidthDL_constr_3, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_carrierBandwidthDL_specs_3	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_nrofCRS_Ports_value2enum_13[] = {
	{ 0,	2,	"n1" },
	{ 1,	2,	"n2" },
	{ 2,	2,	"n4" }
};
static const unsigned int asn_MAP_ASN_RRC_nrofCRS_Ports_enum2value_13[] = {
	0,	/* n1(0) */
	1,	/* n2(1) */
	2	/* n4(2) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_nrofCRS_Ports_specs_13 = {
	asn_MAP_ASN_RRC_nrofCRS_Ports_value2enum_13,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_nrofCRS_Ports_enum2value_13,	/* N => "tag"; sorted by N */
	3,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13[] = {
	(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_nrofCRS_Ports_13 = {
	"nrofCRS-Ports",
	"nrofCRS-Ports",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13,
	sizeof(asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13)
		/sizeof(asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13)
		/sizeof(asn_DEF_ASN_RRC_nrofCRS_Ports_tags_13[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_nrofCRS_Ports_constr_13, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_nrofCRS_Ports_specs_13	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_v_Shift_value2enum_17[] = {
	{ 0,	2,	"n0" },
	{ 1,	2,	"n1" },
	{ 2,	2,	"n2" },
	{ 3,	2,	"n3" },
	{ 4,	2,	"n4" },
	{ 5,	2,	"n5" }
};
static const unsigned int asn_MAP_ASN_RRC_v_Shift_enum2value_17[] = {
	0,	/* n0(0) */
	1,	/* n1(1) */
	2,	/* n2(2) */
	3,	/* n3(3) */
	4,	/* n4(4) */
	5	/* n5(5) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_v_Shift_specs_17 = {
	asn_MAP_ASN_RRC_v_Shift_value2enum_17,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_v_Shift_enum2value_17,	/* N => "tag"; sorted by N */
	6,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_v_Shift_tags_17[] = {
	(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_v_Shift_17 = {
	"v-Shift",
	"v-Shift",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_v_Shift_tags_17,
	sizeof(asn_DEF_ASN_RRC_v_Shift_tags_17)
		/sizeof(asn_DEF_ASN_RRC_v_Shift_tags_17[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_v_Shift_tags_17,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_v_Shift_tags_17)
		/sizeof(asn_DEF_ASN_RRC_v_Shift_tags_17[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_v_Shift_constr_17, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_v_Shift_specs_17	/* Additional specs */
};

asn_TYPE_member_t asn_MBR_ASN_RRC_RateMatchPatternLTE_CRS_1[] = {
	{ ATF_NOFLAGS, 0, offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, carrierFreqDL),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_NativeInteger,
		0,
		{ 0, &asn_PER_memb_ASN_RRC_carrierFreqDL_constr_2,  memb_ASN_RRC_carrierFreqDL_constraint_1 },
		0, 0, /* No default value */
		"carrierFreqDL"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, carrierBandwidthDL),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_carrierBandwidthDL_3,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"carrierBandwidthDL"
		},
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, mbsfn_SubframeConfigList),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_EUTRA_MBSFN_SubframeConfigList,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"mbsfn-SubframeConfigList"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, nrofCRS_Ports),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_nrofCRS_Ports_13,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"nrofCRS-Ports"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, v_Shift),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_v_Shift_17,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"v-Shift"
		},
};
static const int asn_MAP_ASN_RRC_RateMatchPatternLTE_CRS_oms_1[] = { 2 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_RateMatchPatternLTE_CRS_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* carrierFreqDL */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* carrierBandwidthDL */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* mbsfn-SubframeConfigList */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* nrofCRS-Ports */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 } /* v-Shift */
};
asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_RateMatchPatternLTE_CRS_specs_1 = {
	sizeof(struct ASN_RRC_RateMatchPatternLTE_CRS),
	offsetof(struct ASN_RRC_RateMatchPatternLTE_CRS, _asn_ctx),
	asn_MAP_ASN_RRC_RateMatchPatternLTE_CRS_tag2el_1,
	5,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_RateMatchPatternLTE_CRS_oms_1,	/* Optional members */
	1, 0,	/* Root/Additions */
	-1,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS = {
	"RateMatchPatternLTE-CRS",
	"RateMatchPatternLTE-CRS",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1,
	sizeof(asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1)
		/sizeof(asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1[0]), /* 1 */
	asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1)
		/sizeof(asn_DEF_ASN_RRC_RateMatchPatternLTE_CRS_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_RateMatchPatternLTE_CRS_1,
	5,	/* Elements count */
	&asn_SPC_ASN_RRC_RateMatchPatternLTE_CRS_specs_1	/* Additional specs */
};

