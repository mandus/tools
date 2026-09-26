/* sample.c -- exercises Comment, PreProc, Type, Statement, Constant, String,
 * Number, Special, Operator, Delimiter, Todo.
 * TODO: this comment should stand out. FIXME too. XXX as well.
 */
#include <stdio.h>
#include <stdlib.h>
#define MAX_ITEMS 64
#define SQUARE(x) ((x) * (x))

typedef enum { ST_IDLE, ST_BUSY, ST_DONE } state_t;

typedef struct {
    const char *name;
    unsigned    id;
    double      weight;
    state_t     state;
} item_t;

static const char *STATES[] = { "idle", "busy", "done" };

/* A function with every constant flavour. */
int main(int argc, char **argv)
{
    static item_t items[MAX_ITEMS];
    int      count  = 0;
    unsigned mask   = 0xDEADBEEF;
    double   ratio  = 1.618e-3;
    char     tab    = '\t';
    _Bool    ready  = 1;

    if (argc < 2 || argv == NULL) {
        fprintf(stderr, "usage: %s <n>\tescape: \\n \"quoted\"\n", argv[0]);
        return EXIT_FAILURE;
    }

    for (int i = 0; i < MAX_ITEMS && ready; ++i) {
        items[i].name   = STATES[i % 3];
        items[i].id     = (unsigned)i ^ mask;
        items[i].weight = ratio * SQUARE(i);
        items[i].state  = (i & 1) ? ST_BUSY : ST_IDLE;
        count++;
    }

    switch (items[0].state) {
    case ST_IDLE:  puts("idle");  break;
    case ST_BUSY:  puts("busy");  break;
    default:       goto done;
    }

done:
    while (count-- > 0 && tab != '\0')
        continue;

    return EXIT_SUCCESS;
}
