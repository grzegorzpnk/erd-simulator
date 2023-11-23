import logging
import json
import jsonschema
import os
from sb3_contrib import MaskablePPO


def validate_body(data: dict):
    log = logging.getLogger('ermodel')

    schema_path = os.environ.get("GET_PREDICTION_SCHEMA_PATH")

    with open(schema_path) as file:
        schema_json = file.read()

    schema = json.loads(schema_json)

    try:
        jsonschema.validate(data, schema)
        return True
    except jsonschema.exceptions.ValidationError as e:
        log.error("Body schema validation error:", e.message)

        return False


def make_prediction(models: dict, state_t0, mask, use_mask: bool):
    log = logging.getLogger('ermodel')

    if use_mask:
        action, state_t1 = models["masked"].predict(state_t0, action_masks=mask)
    else:
        action, state_t1 = models["basic"].predict(state_t0, deterministic = True)

    # State representation:
    #   APP  : [0,1] Required mvCPU  [0,2] required Memory [0,3] Required Latency [0,4] Current MEC [0,5] Current RAN
    #   MEC  : 0) CPU Capacity 1) CPU Utilization [%] 2) Memory Capacity 3) Memory Utilization [%] 4) Unit Cost
    log.info(f"[INFO] Action[{action}],\nState t0 [{state_t0}],\nState t1[{state_t1}]")

    return action, state_t1


