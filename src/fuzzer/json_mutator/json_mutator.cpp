#include <curl/curl.h>
#include <nlohmann/json.hpp>

#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <random>

using namespace std;
using json = nlohmann::json;

// constants
const int VALUE_TYPES = 4;  // number of valid value types considered
const int MUTATION_NUM = 4; // number of mutation types
const string charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

// random distributions
uniform_int_distribution<> string_length_dist(0, 10);
uniform_int_distribution<> int_dist(numeric_limits<int>::min(), numeric_limits<int>::max());
uniform_real_distribution<> float_dist(numeric_limits<float>::min(), numeric_limits<float>::max());
uniform_int_distribution<> bool_dist(numeric_limits<bool>::min(), numeric_limits<bool>::max());
uniform_int_distribution<> mutate_dist(1, MUTATION_NUM);
uniform_int_distribution<> value_types_dist(1, VALUE_TYPES);
uniform_int_distribution<> char_dist(0, charset.length() - 1);

/* ============================= HELPER FUNCTIONS ============================= */
// generates a random string of size "length"
string generate_random_string(int length)
{
    random_device rd;
    mt19937 gen(rd());

    string random_string;
    for (int i = 0; i < length; ++i)
    {
        random_string += charset[char_dist(gen)];
    }
    return random_string;
}

// checks if a string is convertible to integer
bool string_convertible_to_int(const string &str)
{
    stringstream ss(str);
    int num;
    ss >> num;
    return !ss.fail() && ss.eof();
}

/* ============================= MUTATOR FUNCTIONS ============================= */
// Randomize current value without changing its type
void randomize_value(json &data, string key)
{
    int num_value;
    string string_value;
    float float_value;
    bool bool_value;

    random_device rd;
    mt19937 gen(rd());

    auto value = data[key];

    if (value.is_number())
    {
        num_value = int_dist(gen);
        data[key] = num_value;
    }
    else if (value.is_string())
    {
        string_value = generate_random_string(string_length_dist(gen));
        data[key] = string_value;
    }
    else if (value.is_number_float())
    {
        float_value = float_dist(gen);
        data[key] = float_value;
    }
    else if (value.is_boolean())
    {
        bool_value = bool_dist(gen);
        data[key] = bool_value;
    }
    else if (value.is_null())
    {
        // TODO: find ways to fuzz NULL values
        num_value = int_dist(gen);
        data[key] = num_value;
    }
    else if (value.is_array())
    {
        // TODO: find ways to fuzz arrays
        num_value = int_dist(gen);
        data[key] = num_value;
    }
    else
    {
        cout << "value type is not covered" << endl;
    }
}

// changes one value type to another type
void change_value_type(json &data, const string key)
{
    // random device generators
    random_device rd;
    mt19937 gen(rd());
    int rand;

    // check if key exists
    if (data.find(key) != data.end())
    {
        auto value = data[key];
        json::value_t val_type = data[key].type();

        // 1. UNSIGNED INT or INTEGER
        if (val_type == json::value_t::number_unsigned || val_type == json::value_t::number_integer)
        {
            unsigned int num_val;
            rand = bool_dist(gen);
            if (val_type == json::value_t::number_unsigned)
                num_val = value.get<json::number_unsigned_t>();
            else
                num_val = value.get<json::number_integer_t>();
            if (rand)
            {
                // convert INT to FLOAT
                float float_res = static_cast<float>(num_val);
                data[key] = float_res;
            }
            else
            {
                // convert INT to STRING
                string string_res = to_string(num_val);
                data[key] = string_res;
            }
        }
        // 2. FLOAT
        else if (val_type == json::value_t::number_float)
        {
            // convert FLOAT to INT
            data[key] = static_cast<int>(value);
        }
        // 3. STRING
        else if (val_type == json::value_t::string)
        {
            // convert STRING to INT, BOOL
            json::string_t string_val = value.get<json::string_t>();
            if (string_convertible_to_int(value))
            {
                // convert STRING to INT
                int int_val = stoi(string_val);
                data[key] = int_val;
            }
            else
            {
                // currently just randomize strings containing chars and numbers
                randomize_value(data, key);
            }
        }
        // 4. BOOLEAN
        else if (val_type == json::value_t::boolean)
        {
            // convert BOOLEAN to STRING
            if (data[key])
                data[key] = "true";
            else
                data[key] = "false";
        }
        else
        {
            cout << "No valid type found for this value" << endl;
        }
    }
}

/* ============================= JSON MUTATOR FUNCTION ============================= */
json mutate_requests(string request_type, json &data)
{
    // initialize variables for switch block
    string new_key, new_val;
    int index;
    auto it = data.begin();
    float float_val;

    // create random device generator
    random_device rd;
    mt19937 gen(rd());
    // distribution for random JSON field index
    uniform_int_distribution<> index_dist(0, data.size() - 1);

    if (request_type == "POST")
    {
        int rand = mutate_dist(gen);
        if (data.empty())
        {
            rand = 1; // force to add a new field
        }

        switch (rand)
        {
        // 1. Add new JSON field
        case 1:
            cout << "Mutation 1: Add new JSON field" << endl;
            new_key = generate_random_string(5);
            rand = value_types_dist(gen);
            switch (rand)
            {
            case 1:
                new_val = generate_random_string(5);
                break;
            case 2:
                new_val = to_string(int_dist(gen));
                break;
            case 3:
                new_val = to_string(float_dist(gen));
                break;
            case 4:
                new_val = to_string(bool_dist(gen));
                break;
            }
            data.emplace(new_key, new_val);
            break;

        // 2. Remove existing JSON field
        case 2:
            index = index_dist(gen);
            cout << "Mutation 2: Remove existing JSON field at index " << index << endl;
            it = next(data.begin(), index);
            data.erase(it);
            break;

        // 3. Keep all fields as it is but change value type
        case 3:
            cout << "Mutation 3: Change type of existing JSON field value" << endl;
            index = index_dist(gen);
            it = next(data.begin(), index);
            change_value_type(data, it.key());
            break;

        // 4. Randomize current JSON value
        case 4:
            index = index_dist(gen);
            it = next(data.begin(), index);
            cout << "Mutation 4: Randomize value at index " << index << endl;
            randomize_value(data, it.key());
            break;
        }
    }
    else if (request_type == "GET")
    {
        // TODO: handle GET request JSON mutations
    }

    return data;
}

/* MAIN FUNCTION for TESTING */
/* UNCOMMENT FOR TESTING, as MAIN here clashes with main.cpp */
int main(int argc, char *argv[])
{
    // command line arguments
    int rounds = 10;
    if (argc > 1)
    {
        rounds = stoi(argv[1]);
    }
    cout << "Number of mutation rounds: " << rounds << endl;

    // JSON input file
    string input_file_path = "../src/Django_Web/template_inputs/add_product.json";
    ifstream input_file(input_file_path);
    if (!input_file.is_open())
    {
        cerr << "Error: Could not open file '" << input_file_path << "'." << endl;
        return 1;
    }
    json data;
    input_file >> data;
    input_file.close();

    // loop through "rounds" and mutate in each iteration
    for (int i = 0; i < rounds; i++)
    {
        data = mutate_requests("POST", data);
        if (!data.empty())
            cout << data.dump(4) << endl;
        else
            cout << "Empty data: {}" << endl;
    }

    return 0;
}