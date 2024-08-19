#include <Python.h>
#include <iostream>
#include <vector>
#include <string>

struct Data {
    std::string string_value;
    double numeric_value;
};

void set_path() {
    PyObject *sys = PyImport_ImportModule("sys");
    PyObject *path = PyObject_GetAttrString(sys, "path");
    PyList_Append(path, PyUnicode_DecodeFSDefault("."));
}

int main() {
    Py_Initialize();
    set_path();

    PyObject *pName = PyUnicode_DecodeFSDefault("pytest");

    PyObject *pModule = PyImport_Import(pName);
    Py_DECREF(pName);

    if (pModule != nullptr) {
        PyObject *pFunc = PyObject_GetAttrString(pModule, "get_data");

        if (pFunc && PyCallable_Check(pFunc)) {
            // PyObject *pArgs = PyTuple_New(1);
            // PyObject *pValue = PyUnicode_DecodeFSDefault("Hello from C++");
            // if (!pValue) {
            //     Py_DECREF(pArgs);
            //     Py_DECREF(pModule);
            //     std::cout << "Cannot convert argument" << std::endl;
            //     return 1;
            // }
            // PyTuple_SetItem(pArgs, 0, pValue);
            // Py_DECREF(pValue);

            PyObject *pList = PyObject_CallObject(pFunc, nullptr);
            // Py_DECREF(pArgs);
            
            if (pList != nullptr && PyList_Check(pList)) {
                std::vector<Data> dataList;

                Py_ssize_t listSize = PyList_Size(pList);
                for (Py_ssize_t i = 0; i < listSize; i++) {
                    PyObject *pDict = PyList_GetItem(pList, i);

                    if (PyDict_Check(pDict)) {
                        PyObject *pStringValue = PyDict_GetItemString(pDict, "string_value");
                        PyObject *pNumericValue = PyDict_GetItemString(pDict, "numeric_value");

                        if (PyUnicode_Check(pStringValue) && PyFloat_Check(pNumericValue)) {
                            Data data;
                            data.string_value = PyUnicode_AsUTF8(pStringValue);
                            data.numeric_value = PyFloat_AsDouble(pNumericValue);
                            dataList.push_back(data);
                        }
                        //Py_DECREF(pStringValue);
                        //Py_DECREF(pNumericValue);
                    }
                    //Py_DECREF(pDict);
                }

                for (const auto& data : dataList) {
                    std::cout << "String value: " << data.string_value << ", Numeric value: " << data.numeric_value << std::endl;
                }

                Py_DECREF(pList);
            } else {
                Py_DECREF(pFunc);
                Py_DECREF(pModule);
                std::cout << "Failed to get list" << std::endl;
                PyErr_Print();
                return 1;
            }

        } else {
            if(PyErr_Occurred()) {
                PyErr_Print();
                std::cout << "Failed to get function" << std::endl;
            }
        }
        Py_XDECREF(pFunc);
        Py_DECREF(pModule);
    } else {
        std::cout << "Failed to import module" << std::endl;
        PyErr_Print();
        return 1;
    }

    if (Py_FinalizeEx() < 0) {
        return 120;
    }

    return 0;
}
