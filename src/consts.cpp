#include "consts.h"
#include <map>
#include <string>

const std::map<long, std::string> HTTP_STATUS_MESSAGES = {
    {200, "OK"},
    {201, "Created"},
    {202, "Accepted"},
    {300, "Permanent Redirect"},
    {302, "Temporary Redirect"},
    {304, "Not Modified"},
    {400, "Bad Request"},
    {401, "Unauthorized"},
    {403, "Forbidden"},
    {404, "Not Found"},
    {429, "Too Many Requests"},
    {500, "Internal Server Error"},
    {501, "Not Implemented"},
    {502, "Bad Gateway"},
    {503, "Service Unavailable"},
    {504, "Gateway Timeout"}
};