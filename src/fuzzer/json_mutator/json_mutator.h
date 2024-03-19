#ifndef JSON_MUTATOR_H
#define JSON_MUTATOR_H

#include <nlohmann/json.hpp>

nlohmann::json mutate_requests(std::string request_type, nlohmann::json &data);

#endif