package schema

const promptSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["apiVersion", "kind", "metadata", "spec"],
  "properties": {
    "apiVersion": {"type": "string"},
    "kind": {"type": "string", "const": "Prompt"},
    "metadata": {
      "type": "object",
      "required": ["id"],
      "properties": {
        "id": {"type": "string", "pattern": "^[a-zA-Z0-9_-]+$"},
        "name": {"type": "string"},
        "description": {"type": "string"}
      }
    },
    "spec": {
      "type": "object",
      "required": ["body"],
      "properties": {
        "fragments": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["body"],
            "properties": {
              "inputs": {
                "type": "object",
                "additionalProperties": {
                  "type": "object",
                  "required": ["type"],
                  "properties": {
                    "type": {"type": "string", "enum": ["string", "number", "boolean", "array", "object"]},
                    "required": {"type": "boolean"},
                    "default": {},
                    "items": {"type": "object"},
                    "properties": {"type": "object"}
                  }
                }
              },
              "body": {"type": "string"}
            }
          }
        },
        "body": {
          "oneOf": [
            {"type": "string"},
            {
              "type": "array",
              "items": {
                "oneOf": [
                  {"type": "string"},
                  {
                    "type": "object",
                    "required": ["fragment"],
                    "properties": {
                      "fragment": {"type": "string"},
                      "inputs": {"type": "object"}
                    }
                  }
                ]
              }
            }
          ]
        }
      }
    }
  }
}`

const promptsetSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["apiVersion", "kind", "metadata", "spec"],
  "properties": {
    "apiVersion": {"type": "string"},
    "kind": {"type": "string", "const": "Promptset"},
    "metadata": {
      "type": "object",
      "required": ["id"],
      "properties": {
        "id": {"type": "string", "pattern": "^[a-zA-Z0-9_-]+$"},
        "name": {"type": "string"},
        "description": {"type": "string"}
      }
    },
    "spec": {
      "type": "object",
      "required": ["prompts"],
      "properties": {
        "fragments": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["body"],
            "properties": {
              "inputs": {
                "type": "object",
                "additionalProperties": {
                  "type": "object",
                  "required": ["type"],
                  "properties": {
                    "type": {"type": "string", "enum": ["string", "number", "boolean", "array", "object"]},
                    "required": {"type": "boolean"},
                    "default": {},
                    "items": {"type": "object"},
                    "properties": {"type": "object"}
                  }
                }
              },
              "body": {"type": "string"}
            }
          }
        },
        "prompts": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["body"],
            "properties": {
              "name": {"type": "string"},
              "description": {"type": "string"},
              "body": {
                "oneOf": [
                  {"type": "string"},
                  {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {"type": "string"},
                        {
                          "type": "object",
                          "required": ["fragment"],
                          "properties": {
                            "fragment": {"type": "string"},
                            "inputs": {"type": "object"}
                          }
                        }
                      ]
                    }
                  }
                ]
              }
            }
          }
        }
      }
    }
  }
}`

const ruleSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["apiVersion", "kind", "metadata", "spec"],
  "properties": {
    "apiVersion": {"type": "string"},
    "kind": {"type": "string", "const": "Rule"},
    "metadata": {
      "type": "object",
      "required": ["id"],
      "properties": {
        "id": {"type": "string", "pattern": "^[a-zA-Z0-9_-]+$"},
        "name": {"type": "string"},
        "description": {"type": "string"}
      }
    },
    "spec": {
      "type": "object",
      "required": ["enforcement", "body"],
      "properties": {
        "fragments": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["body"],
            "properties": {
              "inputs": {
                "type": "object",
                "additionalProperties": {
                  "type": "object",
                  "required": ["type"],
                  "properties": {
                    "type": {"type": "string", "enum": ["string", "number", "boolean", "array", "object"]},
                    "required": {"type": "boolean"},
                    "default": {},
                    "items": {"type": "object"},
                    "properties": {"type": "object"}
                  }
                }
              },
              "body": {"type": "string"}
            }
          }
        },
        "enforcement": {"type": "string", "enum": ["may", "should", "must"]},
        "scope": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "files": {
                "type": "array",
                "items": {"type": "string"}
              }
            }
          }
        },
        "body": {
          "oneOf": [
            {"type": "string"},
            {
              "type": "array",
              "items": {
                "oneOf": [
                  {"type": "string"},
                  {
                    "type": "object",
                    "required": ["fragment"],
                    "properties": {
                      "fragment": {"type": "string"},
                      "inputs": {"type": "object"}
                    }
                  }
                ]
              }
            }
          ]
        }
      }
    }
  }
}`

const rulesetSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["apiVersion", "kind", "metadata", "spec"],
  "properties": {
    "apiVersion": {"type": "string"},
    "kind": {"type": "string", "const": "Ruleset"},
    "metadata": {
      "type": "object",
      "required": ["id"],
      "properties": {
        "id": {"type": "string", "pattern": "^[a-zA-Z0-9_-]+$"},
        "name": {"type": "string"},
        "description": {"type": "string"}
      }
    },
    "spec": {
      "type": "object",
      "required": ["rules"],
      "properties": {
        "fragments": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["body"],
            "properties": {
              "inputs": {
                "type": "object",
                "additionalProperties": {
                  "type": "object",
                  "required": ["type"],
                  "properties": {
                    "type": {"type": "string", "enum": ["string", "number", "boolean", "array", "object"]},
                    "required": {"type": "boolean"},
                    "default": {},
                    "items": {"type": "object"},
                    "properties": {"type": "object"}
                  }
                }
              },
              "body": {"type": "string"}
            }
          }
        },
        "rules": {
          "type": "object",
          "additionalProperties": {
            "type": "object",
            "required": ["enforcement", "body"],
            "properties": {
              "name": {"type": "string"},
              "description": {"type": "string"},
              "priority": {"type": "integer"},
              "enforcement": {"type": "string", "enum": ["may", "should", "must"]},
              "scope": {
                "type": "array",
                "items": {
                  "type": "object",
                  "properties": {
                    "files": {
                      "type": "array",
                      "items": {"type": "string"}
                    }
                  }
                }
              },
              "body": {
                "oneOf": [
                  {"type": "string"},
                  {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {"type": "string"},
                        {
                          "type": "object",
                          "required": ["fragment"],
                          "properties": {
                            "fragment": {"type": "string"},
                            "inputs": {"type": "object"}
                          }
                        }
                      ]
                    }
                  }
                ]
              }
            }
          }
        }
      }
    }
  }
}`
