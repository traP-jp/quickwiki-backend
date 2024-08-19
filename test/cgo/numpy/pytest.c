#include </usr/include/python3.10/Python.h>
#include <stdio.h>
#include <string.h>

char result[16];

void set_path() {
    PyObject *sys = PyImport_ImportModule("sys");
    PyObject *path = PyObject_GetAttrString(sys, "path");
    PyList_Append(path, PyUnicode_DecodeFSDefault("."));
}

char *pyHello(char *name, char *fileName, char *funcName) {
    PyObject *pName, *pModule, *pFunc, *pArgs, *pValue;

    Py_Initialize();
    set_path();

    pName = PyUnicode_DecodeFSDefault(fileName);
    pModule = PyImport_Import(pName);
    Py_DECREF(pName);

    if (pModule != NULL) {
        pFunc = PyObject_GetAttrString(pModule, funcName);

        if (pFunc && PyCallable_Check(pFunc)) {
            pArgs = PyTuple_New(1);
            pValue = PyUnicode_DecodeFSDefault(name);
            if (!pValue) {
                Py_DECREF(pArgs);
                Py_DECREF(pModule);
                fprintf(stderr, "Cannot convert argument\n");
                return "Argument Convert Error";
            }
            PyTuple_SetItem(pArgs, 0, pValue);

            pValue = PyObject_CallObject(pFunc, pArgs);
            Py_DECREF(pArgs);
            if (pValue != NULL) {
                printf("Result of call: %s\n", PyUnicode_AsUTF8(pValue));
                strcpy(result, PyUnicode_AsUTF8(pValue));
                Py_DECREF(pValue);
            } else {
                Py_DECREF(pFunc);
                Py_DECREF(pModule);
                PyErr_Print();
                fprintf(stderr, "Function call failed\n");
                return "Function Call Error";
            }
        } else {
            if (PyErr_Occurred()) {
                PyErr_Print();
            }
            fprintf(stderr, "Cannot find function \"%s\"\n", funcName);
        }
        Py_XDECREF(pFunc);
        Py_DECREF(pModule);
    } else {
        PyErr_Print();
        fprintf(stderr, "Failed to load \"%s.py\"\n", fileName);
        return "File Loading Error";
    }

    if (Py_FinalizeEx() < 0) {
        return "Finalize Error";
    }

    printf("Finish calling python\nresult: %s\n", result);
    return result;
}

int main() {
    pyHello("World", "hello", "hello");
}