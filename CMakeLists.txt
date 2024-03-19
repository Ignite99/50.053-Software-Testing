cmake_minimum_required(VERSION 3.10)

set(CMAKE_CXX_COMPILER g++)
project(Software_Testing_Project)

# Set C++ standard to C++17
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Add the executable target
add_executable(Software_Testing_Project
    src/main.cpp
    src/fuzzer/fuzzer.cpp
    src/BLE_Zephyr/ble_zephyr.cpp
    src/CoAP_Protocol/coap_protocol.cpp
    src/Django_Web/django_web.cpp
    src/HTML_Logger/html_logger.cpp
)

# Include directories
target_include_directories(Software_Testing_Project PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/nlohmann
)

find_package(CURL REQUIRED)
target_link_libraries(Software_Testing_Project PRIVATE ${CURL_LIBRARIES})

# Add consts.cpp to the target
target_sources(Software_Testing_Project PRIVATE src/consts.cpp)