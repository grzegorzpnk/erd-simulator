import json
import logging


def default_config_file():
    config = json.loads(
        '''
        {
            "model_path": "./model"
        }
        '''
    )
    return config


class Config:
    def __init__(self, config_path):
        self.logger = logging.getLogger("django")
        self.config = dict(self.read_config_file(config_path))

    def get_model_path(self):
        return self.config['model_path']

    def read_config_file(self, path):
        try:
            with open(path) as cfg:
                config = json.load(cfg)
        except FileExistsError:
            self.logger.error(f"Error: Config file not found, path: {path}. Using default config: {default_config_file()}")
            return default_config_file()
        except json.JSONDecodeError:
            self.logger.error(f"Error: Invalid JSON syntax in config file. Using default config: {default_config_file()}")
            return default_config_file()
        except Exception as e:
            self.logger.error(f"Error: {str(e)}. Using default config: {default_config_file()}")
            return default_config_file()

        return config
