/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "NR-RRC-Definitions"
 * 	found in "asn/nr-rrc-15.6.0.asn1"
 * 	`asn1c -fcompound-names -pdu=all -findirect-choice -fno-include-deps -gen-PER -no-gen-OER -no-gen-example -D rrc`
 */

#include "ASN_RRC_MeasAndMobParametersFRX-Diff.h"

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
/*
 * This type is implemented using NativeEnumerated,
 * so here we adjust the DEF accordingly.
 */
static asn_per_constraints_t asn_PER_type_ASN_RRC_ss_SINR_Meas_constr_2 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_constr_4 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_constr_6 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_csi_SINR_Meas_constr_8 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_csi_RS_RLM_constr_10 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_handoverInterF_constr_14 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_handoverLTE_EPC_constr_16 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_handoverLTE_5GC_constr_18 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_maxNumberResource_CSI_RS_RLM_constr_21 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 2,  2,  0,  3 }	/* (0..3) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_constr_27 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_ss_SINR_Meas_value2enum_2[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_ss_SINR_Meas_enum2value_2[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_ss_SINR_Meas_specs_2 = {
	asn_MAP_ASN_RRC_ss_SINR_Meas_value2enum_2,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_ss_SINR_Meas_enum2value_2,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2[] = {
	(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_ss_SINR_Meas_2 = {
	"ss-SINR-Meas",
	"ss-SINR-Meas",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2,
	sizeof(asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2)
		/sizeof(asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2)
		/sizeof(asn_DEF_ASN_RRC_ss_SINR_Meas_tags_2[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_ss_SINR_Meas_constr_2, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_ss_SINR_Meas_specs_2	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_value2enum_4[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_enum2value_4[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_specs_4 = {
	asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_value2enum_4,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_enum2value_4,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4[] = {
	(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_4 = {
	"csi-RSRP-AndRSRQ-MeasWithSSB",
	"csi-RSRP-AndRSRQ-MeasWithSSB",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4,
	sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4)
		/sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4)
		/sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_tags_4[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_constr_4, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_specs_4	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_value2enum_6[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_enum2value_6[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_specs_6 = {
	asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_value2enum_6,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_enum2value_6,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6[] = {
	(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_6 = {
	"csi-RSRP-AndRSRQ-MeasWithoutSSB",
	"csi-RSRP-AndRSRQ-MeasWithoutSSB",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6,
	sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6)
		/sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6)
		/sizeof(asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_tags_6[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_constr_6, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_specs_6	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_csi_SINR_Meas_value2enum_8[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_csi_SINR_Meas_enum2value_8[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_csi_SINR_Meas_specs_8 = {
	asn_MAP_ASN_RRC_csi_SINR_Meas_value2enum_8,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_csi_SINR_Meas_enum2value_8,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8[] = {
	(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_csi_SINR_Meas_8 = {
	"csi-SINR-Meas",
	"csi-SINR-Meas",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8,
	sizeof(asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8)
		/sizeof(asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8)
		/sizeof(asn_DEF_ASN_RRC_csi_SINR_Meas_tags_8[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_csi_SINR_Meas_constr_8, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_csi_SINR_Meas_specs_8	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_csi_RS_RLM_value2enum_10[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_csi_RS_RLM_enum2value_10[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_csi_RS_RLM_specs_10 = {
	asn_MAP_ASN_RRC_csi_RS_RLM_value2enum_10,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_csi_RS_RLM_enum2value_10,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_csi_RS_RLM_tags_10[] = {
	(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_csi_RS_RLM_10 = {
	"csi-RS-RLM",
	"csi-RS-RLM",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_csi_RS_RLM_tags_10,
	sizeof(asn_DEF_ASN_RRC_csi_RS_RLM_tags_10)
		/sizeof(asn_DEF_ASN_RRC_csi_RS_RLM_tags_10[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_csi_RS_RLM_tags_10,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_csi_RS_RLM_tags_10)
		/sizeof(asn_DEF_ASN_RRC_csi_RS_RLM_tags_10[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_csi_RS_RLM_constr_10, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_csi_RS_RLM_specs_10	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_handoverInterF_value2enum_14[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_handoverInterF_enum2value_14[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_handoverInterF_specs_14 = {
	asn_MAP_ASN_RRC_handoverInterF_value2enum_14,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_handoverInterF_enum2value_14,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_handoverInterF_tags_14[] = {
	(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_handoverInterF_14 = {
	"handoverInterF",
	"handoverInterF",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_handoverInterF_tags_14,
	sizeof(asn_DEF_ASN_RRC_handoverInterF_tags_14)
		/sizeof(asn_DEF_ASN_RRC_handoverInterF_tags_14[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_handoverInterF_tags_14,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_handoverInterF_tags_14)
		/sizeof(asn_DEF_ASN_RRC_handoverInterF_tags_14[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_handoverInterF_constr_14, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_handoverInterF_specs_14	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_handoverLTE_EPC_value2enum_16[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_handoverLTE_EPC_enum2value_16[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_handoverLTE_EPC_specs_16 = {
	asn_MAP_ASN_RRC_handoverLTE_EPC_value2enum_16,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_handoverLTE_EPC_enum2value_16,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16[] = {
	(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_handoverLTE_EPC_16 = {
	"handoverLTE-EPC",
	"handoverLTE-EPC",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16,
	sizeof(asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16)
		/sizeof(asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16)
		/sizeof(asn_DEF_ASN_RRC_handoverLTE_EPC_tags_16[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_handoverLTE_EPC_constr_16, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_handoverLTE_EPC_specs_16	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_handoverLTE_5GC_value2enum_18[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_handoverLTE_5GC_enum2value_18[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_handoverLTE_5GC_specs_18 = {
	asn_MAP_ASN_RRC_handoverLTE_5GC_value2enum_18,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_handoverLTE_5GC_enum2value_18,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18[] = {
	(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_handoverLTE_5GC_18 = {
	"handoverLTE-5GC",
	"handoverLTE-5GC",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18,
	sizeof(asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18)
		/sizeof(asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18)
		/sizeof(asn_DEF_ASN_RRC_handoverLTE_5GC_tags_18[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_handoverLTE_5GC_constr_18, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_handoverLTE_5GC_specs_18	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_ext1_13[] = {
	{ ATF_POINTER, 3, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext1, handoverInterF),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_handoverInterF_14,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"handoverInterF"
		},
	{ ATF_POINTER, 2, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext1, handoverLTE_EPC),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_handoverLTE_EPC_16,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"handoverLTE-EPC"
		},
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext1, handoverLTE_5GC),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_handoverLTE_5GC_18,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"handoverLTE-5GC"
		},
};
static const int asn_MAP_ASN_RRC_ext1_oms_13[] = { 0, 1, 2 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_ext1_tags_13[] = {
	(ASN_TAG_CLASS_CONTEXT | (5 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_ext1_tag2el_13[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* handoverInterF */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* handoverLTE-EPC */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 } /* handoverLTE-5GC */
};
static asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_ext1_specs_13 = {
	sizeof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext1),
	offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext1, _asn_ctx),
	asn_MAP_ASN_RRC_ext1_tag2el_13,
	3,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_ext1_oms_13,	/* Optional members */
	3, 0,	/* Root/Additions */
	-1,	/* First extension addition */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_ext1_13 = {
	"ext1",
	"ext1",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_ext1_tags_13,
	sizeof(asn_DEF_ASN_RRC_ext1_tags_13)
		/sizeof(asn_DEF_ASN_RRC_ext1_tags_13[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_ext1_tags_13,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_ext1_tags_13)
		/sizeof(asn_DEF_ASN_RRC_ext1_tags_13[0]), /* 2 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_ext1_13,
	3,	/* Elements count */
	&asn_SPC_ASN_RRC_ext1_specs_13	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_maxNumberResource_CSI_RS_RLM_value2enum_21[] = {
	{ 0,	2,	"n2" },
	{ 1,	2,	"n4" },
	{ 2,	2,	"n6" },
	{ 3,	2,	"n8" }
};
static const unsigned int asn_MAP_ASN_RRC_maxNumberResource_CSI_RS_RLM_enum2value_21[] = {
	0,	/* n2(0) */
	1,	/* n4(1) */
	2,	/* n6(2) */
	3	/* n8(3) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_maxNumberResource_CSI_RS_RLM_specs_21 = {
	asn_MAP_ASN_RRC_maxNumberResource_CSI_RS_RLM_value2enum_21,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_maxNumberResource_CSI_RS_RLM_enum2value_21,	/* N => "tag"; sorted by N */
	4,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21[] = {
	(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_21 = {
	"maxNumberResource-CSI-RS-RLM",
	"maxNumberResource-CSI-RS-RLM",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21,
	sizeof(asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21)
		/sizeof(asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21)
		/sizeof(asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_tags_21[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_maxNumberResource_CSI_RS_RLM_constr_21, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_maxNumberResource_CSI_RS_RLM_specs_21	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_ext2_20[] = {
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext2, maxNumberResource_CSI_RS_RLM),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_maxNumberResource_CSI_RS_RLM_21,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"maxNumberResource-CSI-RS-RLM"
		},
};
static const int asn_MAP_ASN_RRC_ext2_oms_20[] = { 0 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_ext2_tags_20[] = {
	(ASN_TAG_CLASS_CONTEXT | (6 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_ext2_tag2el_20[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 } /* maxNumberResource-CSI-RS-RLM */
};
static asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_ext2_specs_20 = {
	sizeof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext2),
	offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext2, _asn_ctx),
	asn_MAP_ASN_RRC_ext2_tag2el_20,
	1,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_ext2_oms_20,	/* Optional members */
	1, 0,	/* Root/Additions */
	-1,	/* First extension addition */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_ext2_20 = {
	"ext2",
	"ext2",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_ext2_tags_20,
	sizeof(asn_DEF_ASN_RRC_ext2_tags_20)
		/sizeof(asn_DEF_ASN_RRC_ext2_tags_20[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_ext2_tags_20,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_ext2_tags_20)
		/sizeof(asn_DEF_ASN_RRC_ext2_tags_20[0]), /* 2 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_ext2_20,
	1,	/* Elements count */
	&asn_SPC_ASN_RRC_ext2_specs_20	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_value2enum_27[] = {
	{ 0,	9,	"supported" }
};
static const unsigned int asn_MAP_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_enum2value_27[] = {
	0	/* supported(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_specs_27 = {
	asn_MAP_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_value2enum_27,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_enum2value_27,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27[] = {
	(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_27 = {
	"simultaneousRxDataSSB-DiffNumerology",
	"simultaneousRxDataSSB-DiffNumerology",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27,
	sizeof(asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27)
		/sizeof(asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27)
		/sizeof(asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_tags_27[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_constr_27, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_specs_27	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_ext3_26[] = {
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext3, simultaneousRxDataSSB_DiffNumerology),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_simultaneousRxDataSSB_DiffNumerology_27,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"simultaneousRxDataSSB-DiffNumerology"
		},
};
static const int asn_MAP_ASN_RRC_ext3_oms_26[] = { 0 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_ext3_tags_26[] = {
	(ASN_TAG_CLASS_CONTEXT | (7 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_ext3_tag2el_26[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 } /* simultaneousRxDataSSB-DiffNumerology */
};
static asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_ext3_specs_26 = {
	sizeof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext3),
	offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff__ext3, _asn_ctx),
	asn_MAP_ASN_RRC_ext3_tag2el_26,
	1,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_ext3_oms_26,	/* Optional members */
	1, 0,	/* Root/Additions */
	-1,	/* First extension addition */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_ext3_26 = {
	"ext3",
	"ext3",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_ext3_tags_26,
	sizeof(asn_DEF_ASN_RRC_ext3_tags_26)
		/sizeof(asn_DEF_ASN_RRC_ext3_tags_26[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_ext3_tags_26,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_ext3_tags_26)
		/sizeof(asn_DEF_ASN_RRC_ext3_tags_26[0]), /* 2 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_ext3_26,
	1,	/* Elements count */
	&asn_SPC_ASN_RRC_ext3_specs_26	/* Additional specs */
};

asn_TYPE_member_t asn_MBR_ASN_RRC_MeasAndMobParametersFRX_Diff_1[] = {
	{ ATF_POINTER, 8, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, ss_SINR_Meas),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_ss_SINR_Meas_2,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"ss-SINR-Meas"
		},
	{ ATF_POINTER, 7, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, csi_RSRP_AndRSRQ_MeasWithSSB),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithSSB_4,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"csi-RSRP-AndRSRQ-MeasWithSSB"
		},
	{ ATF_POINTER, 6, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, csi_RSRP_AndRSRQ_MeasWithoutSSB),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_csi_RSRP_AndRSRQ_MeasWithoutSSB_6,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"csi-RSRP-AndRSRQ-MeasWithoutSSB"
		},
	{ ATF_POINTER, 5, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, csi_SINR_Meas),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_csi_SINR_Meas_8,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"csi-SINR-Meas"
		},
	{ ATF_POINTER, 4, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, csi_RS_RLM),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_csi_RS_RLM_10,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"csi-RS-RLM"
		},
	{ ATF_POINTER, 3, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, ext1),
		(ASN_TAG_CLASS_CONTEXT | (5 << 2)),
		0,
		&asn_DEF_ASN_RRC_ext1_13,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"ext1"
		},
	{ ATF_POINTER, 2, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, ext2),
		(ASN_TAG_CLASS_CONTEXT | (6 << 2)),
		0,
		&asn_DEF_ASN_RRC_ext2_20,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"ext2"
		},
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, ext3),
		(ASN_TAG_CLASS_CONTEXT | (7 << 2)),
		0,
		&asn_DEF_ASN_RRC_ext3_26,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"ext3"
		},
};
static const int asn_MAP_ASN_RRC_MeasAndMobParametersFRX_Diff_oms_1[] = { 0, 1, 2, 3, 4, 5, 6, 7 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_MeasAndMobParametersFRX_Diff_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* ss-SINR-Meas */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* csi-RSRP-AndRSRQ-MeasWithSSB */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* csi-RSRP-AndRSRQ-MeasWithoutSSB */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* csi-SINR-Meas */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 }, /* csi-RS-RLM */
    { (ASN_TAG_CLASS_CONTEXT | (5 << 2)), 5, 0, 0 }, /* ext1 */
    { (ASN_TAG_CLASS_CONTEXT | (6 << 2)), 6, 0, 0 }, /* ext2 */
    { (ASN_TAG_CLASS_CONTEXT | (7 << 2)), 7, 0, 0 } /* ext3 */
};
asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_MeasAndMobParametersFRX_Diff_specs_1 = {
	sizeof(struct ASN_RRC_MeasAndMobParametersFRX_Diff),
	offsetof(struct ASN_RRC_MeasAndMobParametersFRX_Diff, _asn_ctx),
	asn_MAP_ASN_RRC_MeasAndMobParametersFRX_Diff_tag2el_1,
	8,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_MeasAndMobParametersFRX_Diff_oms_1,	/* Optional members */
	5, 3,	/* Root/Additions */
	5,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff = {
	"MeasAndMobParametersFRX-Diff",
	"MeasAndMobParametersFRX-Diff",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1,
	sizeof(asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1)
		/sizeof(asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1[0]), /* 1 */
	asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1)
		/sizeof(asn_DEF_ASN_RRC_MeasAndMobParametersFRX_Diff_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_MeasAndMobParametersFRX_Diff_1,
	8,	/* Elements count */
	&asn_SPC_ASN_RRC_MeasAndMobParametersFRX_Diff_specs_1	/* Additional specs */
};

