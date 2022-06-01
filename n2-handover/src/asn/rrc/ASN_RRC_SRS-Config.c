/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "NR-RRC-Definitions"
 * 	found in "asn/nr-rrc-15.6.0.asn1"
 * 	`asn1c -fcompound-names -pdu=all -findirect-choice -fno-include-deps -gen-PER -no-gen-OER -no-gen-example -D rrc`
 */

#include "ASN_RRC_SRS-Config.h"

#include "ASN_RRC_SRS-ResourceSet.h"
#include "ASN_RRC_SRS-Resource.h"
/*
 * This type is implemented using NativeEnumerated,
 * so here we adjust the DEF accordingly.
 */
static int
memb_ASN_RRC_srs_ResourceSetToReleaseList_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	size_t size;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	/* Determine the number of elements */
	size = _A_CSEQUENCE_FROM_VOID(sptr)->count;
	
	if((size >= 1 && size <= 16)) {
		/* Perform validation of the inner elements */
		return td->encoding_constraints.general_constraints(td, sptr, ctfailcb, app_key);
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static int
memb_ASN_RRC_srs_ResourceSetToAddModList_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	size_t size;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	/* Determine the number of elements */
	size = _A_CSEQUENCE_FROM_VOID(sptr)->count;
	
	if((size >= 1 && size <= 16)) {
		/* Perform validation of the inner elements */
		return td->encoding_constraints.general_constraints(td, sptr, ctfailcb, app_key);
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static int
memb_ASN_RRC_srs_ResourceToReleaseList_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	size_t size;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	/* Determine the number of elements */
	size = _A_CSEQUENCE_FROM_VOID(sptr)->count;
	
	if((size >= 1 && size <= 64)) {
		/* Perform validation of the inner elements */
		return td->encoding_constraints.general_constraints(td, sptr, ctfailcb, app_key);
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static int
memb_ASN_RRC_srs_ResourceToAddModList_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	size_t size;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	/* Determine the number of elements */
	size = _A_CSEQUENCE_FROM_VOID(sptr)->count;
	
	if((size >= 1 && size <= 64)) {
		/* Perform validation of the inner elements */
		return td->encoding_constraints.general_constraints(td, sptr, ctfailcb, app_key);
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static asn_per_constraints_t asn_PER_type_ASN_RRC_srs_ResourceSetToReleaseList_constr_2 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 4,  4,  1,  16 }	/* (SIZE(1..16)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_srs_ResourceSetToAddModList_constr_4 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 4,  4,  1,  16 }	/* (SIZE(1..16)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_srs_ResourceToReleaseList_constr_6 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 6,  6,  1,  64 }	/* (SIZE(1..64)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_srs_ResourceToAddModList_constr_8 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 6,  6,  1,  64 }	/* (SIZE(1..64)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_type_ASN_RRC_tpc_Accumulation_constr_10 CC_NOTUSED = {
	{ APC_CONSTRAINED,	 0,  0,  0,  0 }	/* (0..0) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_memb_ASN_RRC_srs_ResourceSetToReleaseList_constr_2 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 4,  4,  1,  16 }	/* (SIZE(1..16)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_memb_ASN_RRC_srs_ResourceSetToAddModList_constr_4 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 4,  4,  1,  16 }	/* (SIZE(1..16)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_memb_ASN_RRC_srs_ResourceToReleaseList_constr_6 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 6,  6,  1,  64 }	/* (SIZE(1..64)) */,
	0, 0	/* No PER value map */
};
static asn_per_constraints_t asn_PER_memb_ASN_RRC_srs_ResourceToAddModList_constr_8 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 6,  6,  1,  64 }	/* (SIZE(1..64)) */,
	0, 0	/* No PER value map */
};
static asn_TYPE_member_t asn_MBR_ASN_RRC_srs_ResourceSetToReleaseList_2[] = {
	{ ATF_POINTER, 0, 0,
		(ASN_TAG_CLASS_UNIVERSAL | (2 << 2)),
		0,
		&asn_DEF_ASN_RRC_SRS_ResourceSetId,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		""
		},
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2[] = {
	(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static asn_SET_OF_specifics_t asn_SPC_ASN_RRC_srs_ResourceSetToReleaseList_specs_2 = {
	sizeof(struct ASN_RRC_SRS_Config__srs_ResourceSetToReleaseList),
	offsetof(struct ASN_RRC_SRS_Config__srs_ResourceSetToReleaseList, _asn_ctx),
	0,	/* XER encoding is XMLDelimitedItemList */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_2 = {
	"srs-ResourceSetToReleaseList",
	"srs-ResourceSetToReleaseList",
	&asn_OP_SEQUENCE_OF,
	asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2,
	sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_tags_2[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_srs_ResourceSetToReleaseList_constr_2, SEQUENCE_OF_constraint },
	asn_MBR_ASN_RRC_srs_ResourceSetToReleaseList_2,
	1,	/* Single element */
	&asn_SPC_ASN_RRC_srs_ResourceSetToReleaseList_specs_2	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_srs_ResourceSetToAddModList_4[] = {
	{ ATF_POINTER, 0, 0,
		(ASN_TAG_CLASS_UNIVERSAL | (16 << 2)),
		0,
		&asn_DEF_ASN_RRC_SRS_ResourceSet,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		""
		},
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4[] = {
	(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static asn_SET_OF_specifics_t asn_SPC_ASN_RRC_srs_ResourceSetToAddModList_specs_4 = {
	sizeof(struct ASN_RRC_SRS_Config__srs_ResourceSetToAddModList),
	offsetof(struct ASN_RRC_SRS_Config__srs_ResourceSetToAddModList, _asn_ctx),
	0,	/* XER encoding is XMLDelimitedItemList */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_4 = {
	"srs-ResourceSetToAddModList",
	"srs-ResourceSetToAddModList",
	&asn_OP_SEQUENCE_OF,
	asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4,
	sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_tags_4[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_srs_ResourceSetToAddModList_constr_4, SEQUENCE_OF_constraint },
	asn_MBR_ASN_RRC_srs_ResourceSetToAddModList_4,
	1,	/* Single element */
	&asn_SPC_ASN_RRC_srs_ResourceSetToAddModList_specs_4	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_srs_ResourceToReleaseList_6[] = {
	{ ATF_POINTER, 0, 0,
		(ASN_TAG_CLASS_UNIVERSAL | (2 << 2)),
		0,
		&asn_DEF_ASN_RRC_SRS_ResourceId,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		""
		},
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6[] = {
	(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static asn_SET_OF_specifics_t asn_SPC_ASN_RRC_srs_ResourceToReleaseList_specs_6 = {
	sizeof(struct ASN_RRC_SRS_Config__srs_ResourceToReleaseList),
	offsetof(struct ASN_RRC_SRS_Config__srs_ResourceToReleaseList, _asn_ctx),
	0,	/* XER encoding is XMLDelimitedItemList */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_srs_ResourceToReleaseList_6 = {
	"srs-ResourceToReleaseList",
	"srs-ResourceToReleaseList",
	&asn_OP_SEQUENCE_OF,
	asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6,
	sizeof(asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceToReleaseList_tags_6[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_srs_ResourceToReleaseList_constr_6, SEQUENCE_OF_constraint },
	asn_MBR_ASN_RRC_srs_ResourceToReleaseList_6,
	1,	/* Single element */
	&asn_SPC_ASN_RRC_srs_ResourceToReleaseList_specs_6	/* Additional specs */
};

static asn_TYPE_member_t asn_MBR_ASN_RRC_srs_ResourceToAddModList_8[] = {
	{ ATF_POINTER, 0, 0,
		(ASN_TAG_CLASS_UNIVERSAL | (16 << 2)),
		0,
		&asn_DEF_ASN_RRC_SRS_Resource,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		""
		},
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8[] = {
	(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static asn_SET_OF_specifics_t asn_SPC_ASN_RRC_srs_ResourceToAddModList_specs_8 = {
	sizeof(struct ASN_RRC_SRS_Config__srs_ResourceToAddModList),
	offsetof(struct ASN_RRC_SRS_Config__srs_ResourceToAddModList, _asn_ctx),
	0,	/* XER encoding is XMLDelimitedItemList */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_srs_ResourceToAddModList_8 = {
	"srs-ResourceToAddModList",
	"srs-ResourceToAddModList",
	&asn_OP_SEQUENCE_OF,
	asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8,
	sizeof(asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8)
		/sizeof(asn_DEF_ASN_RRC_srs_ResourceToAddModList_tags_8[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_srs_ResourceToAddModList_constr_8, SEQUENCE_OF_constraint },
	asn_MBR_ASN_RRC_srs_ResourceToAddModList_8,
	1,	/* Single element */
	&asn_SPC_ASN_RRC_srs_ResourceToAddModList_specs_8	/* Additional specs */
};

static const asn_INTEGER_enum_map_t asn_MAP_ASN_RRC_tpc_Accumulation_value2enum_10[] = {
	{ 0,	8,	"disabled" }
};
static const unsigned int asn_MAP_ASN_RRC_tpc_Accumulation_enum2value_10[] = {
	0	/* disabled(0) */
};
static const asn_INTEGER_specifics_t asn_SPC_ASN_RRC_tpc_Accumulation_specs_10 = {
	asn_MAP_ASN_RRC_tpc_Accumulation_value2enum_10,	/* "tag" => N; sorted by tag */
	asn_MAP_ASN_RRC_tpc_Accumulation_enum2value_10,	/* N => "tag"; sorted by N */
	1,	/* Number of elements in the maps */
	0,	/* Enumeration is not extensible */
	1,	/* Strict enumeration */
	0,	/* Native long size */
	0
};
static const ber_tlv_tag_t asn_DEF_ASN_RRC_tpc_Accumulation_tags_10[] = {
	(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (10 << 2))
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_tpc_Accumulation_10 = {
	"tpc-Accumulation",
	"tpc-Accumulation",
	&asn_OP_NativeEnumerated,
	asn_DEF_ASN_RRC_tpc_Accumulation_tags_10,
	sizeof(asn_DEF_ASN_RRC_tpc_Accumulation_tags_10)
		/sizeof(asn_DEF_ASN_RRC_tpc_Accumulation_tags_10[0]) - 1, /* 1 */
	asn_DEF_ASN_RRC_tpc_Accumulation_tags_10,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_tpc_Accumulation_tags_10)
		/sizeof(asn_DEF_ASN_RRC_tpc_Accumulation_tags_10[0]), /* 2 */
	{ 0, &asn_PER_type_ASN_RRC_tpc_Accumulation_constr_10, NativeEnumerated_constraint },
	0, 0,	/* Defined elsewhere */
	&asn_SPC_ASN_RRC_tpc_Accumulation_specs_10	/* Additional specs */
};

asn_TYPE_member_t asn_MBR_ASN_RRC_SRS_Config_1[] = {
	{ ATF_POINTER, 5, offsetof(struct ASN_RRC_SRS_Config, srs_ResourceSetToReleaseList),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		0,
		&asn_DEF_ASN_RRC_srs_ResourceSetToReleaseList_2,
		0,
		{ 0, &asn_PER_memb_ASN_RRC_srs_ResourceSetToReleaseList_constr_2,  memb_ASN_RRC_srs_ResourceSetToReleaseList_constraint_1 },
		0, 0, /* No default value */
		"srs-ResourceSetToReleaseList"
		},
	{ ATF_POINTER, 4, offsetof(struct ASN_RRC_SRS_Config, srs_ResourceSetToAddModList),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		0,
		&asn_DEF_ASN_RRC_srs_ResourceSetToAddModList_4,
		0,
		{ 0, &asn_PER_memb_ASN_RRC_srs_ResourceSetToAddModList_constr_4,  memb_ASN_RRC_srs_ResourceSetToAddModList_constraint_1 },
		0, 0, /* No default value */
		"srs-ResourceSetToAddModList"
		},
	{ ATF_POINTER, 3, offsetof(struct ASN_RRC_SRS_Config, srs_ResourceToReleaseList),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		0,
		&asn_DEF_ASN_RRC_srs_ResourceToReleaseList_6,
		0,
		{ 0, &asn_PER_memb_ASN_RRC_srs_ResourceToReleaseList_constr_6,  memb_ASN_RRC_srs_ResourceToReleaseList_constraint_1 },
		0, 0, /* No default value */
		"srs-ResourceToReleaseList"
		},
	{ ATF_POINTER, 2, offsetof(struct ASN_RRC_SRS_Config, srs_ResourceToAddModList),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		0,
		&asn_DEF_ASN_RRC_srs_ResourceToAddModList_8,
		0,
		{ 0, &asn_PER_memb_ASN_RRC_srs_ResourceToAddModList_constr_8,  memb_ASN_RRC_srs_ResourceToAddModList_constraint_1 },
		0, 0, /* No default value */
		"srs-ResourceToAddModList"
		},
	{ ATF_POINTER, 1, offsetof(struct ASN_RRC_SRS_Config, tpc_Accumulation),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ASN_RRC_tpc_Accumulation_10,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"tpc-Accumulation"
		},
};
static const int asn_MAP_ASN_RRC_SRS_Config_oms_1[] = { 0, 1, 2, 3, 4 };
static const ber_tlv_tag_t asn_DEF_ASN_RRC_SRS_Config_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ASN_RRC_SRS_Config_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* srs-ResourceSetToReleaseList */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* srs-ResourceSetToAddModList */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* srs-ResourceToReleaseList */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* srs-ResourceToAddModList */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 } /* tpc-Accumulation */
};
asn_SEQUENCE_specifics_t asn_SPC_ASN_RRC_SRS_Config_specs_1 = {
	sizeof(struct ASN_RRC_SRS_Config),
	offsetof(struct ASN_RRC_SRS_Config, _asn_ctx),
	asn_MAP_ASN_RRC_SRS_Config_tag2el_1,
	5,	/* Count of tags in the map */
	asn_MAP_ASN_RRC_SRS_Config_oms_1,	/* Optional members */
	5, 0,	/* Root/Additions */
	5,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_ASN_RRC_SRS_Config = {
	"SRS-Config",
	"SRS-Config",
	&asn_OP_SEQUENCE,
	asn_DEF_ASN_RRC_SRS_Config_tags_1,
	sizeof(asn_DEF_ASN_RRC_SRS_Config_tags_1)
		/sizeof(asn_DEF_ASN_RRC_SRS_Config_tags_1[0]), /* 1 */
	asn_DEF_ASN_RRC_SRS_Config_tags_1,	/* Same as above */
	sizeof(asn_DEF_ASN_RRC_SRS_Config_tags_1)
		/sizeof(asn_DEF_ASN_RRC_SRS_Config_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ASN_RRC_SRS_Config_1,
	5,	/* Elements count */
	&asn_SPC_ASN_RRC_SRS_Config_specs_1	/* Additional specs */
};

