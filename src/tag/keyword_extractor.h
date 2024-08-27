#ifdef __cplusplus
extern "C" {
#endif

#include <Python.h>

void initialize_python();
void finalize_python();
const char* extract(const char *text, int num_keywords);

#ifdef __cplusplus
}
#endif