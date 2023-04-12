import logging
from django.http import HttpRequest, JsonResponse, HttpResponse
from sb3_contrib import MaskablePPO
import os
import sys
from config import config
import json
import jsonschema

sys.path.append("../config")


# Create your views here.
def get_prediction(request: HttpRequest):
    log = logging.getLogger('ermodel')
    if request.method == "GET" or request.method == "POST":
        body = request.body
        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            log.error({'error': 'Invalid JSON'})
            return JsonResponse({'error': 'Invalid JSON'}, status=400)

        if not validate_body(data):
            log.error({'error': 'Body schema validation error'})
            return JsonResponse({'error': 'Body schema validation error'}, status=400)

        state = data['state']
        if "mask" in data:
            use_mask = True
            mask = data['mask']
        else:
            use_mask = False

        log.info(f"Mask present: {use_mask}")

        cfg = config.Config(os.environ.get('CONFIG_PATH'))
        model = MaskablePPO.load(cfg.get_model_path())

        try:
            if use_mask:
                action, _ = model.predict(state, action_masks=mask)
            else:
                action, _ = model.predict(state)
        except ValueError as e:
            log.error({"error": {'state': f'{state}', 'mask': f'{mask}', "exception": f'{e}'}})
            return JsonResponse({"error": {'state': f'{state}', 'mask': f'{mask}', "exception": f'{e}'}}, status=400)

        action += 1

        return HttpResponse(f'{action}')
    else:
        log.error({'error': f'Method [{request.method}] Not Allowed'})
        return JsonResponse({'error': f'Method [{request.method}] Not Allowed'}, status=405)


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
