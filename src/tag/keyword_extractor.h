#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <Python.h>

typedef struct {
    char** tag_names;
    double* scores;
    size_t size;
} DataArray;

void initialize_python();
void finalize_python();
DataArray extract(const char *text, int num_keywords);
void free_data_array(DataArray dataArray);

#ifdef __cplusplus
}
#endif