/* sleep - minimal GNU-sleep-alike, primarily for Windows.
 *
 * Usage: sleep NUMBER[smhd]...
 *   Sleeps for the sum of all arguments. NUMBER may be fractional.
 *   Suffixes: s = seconds (default), m = minutes, h = hours, d = days.
 *
 * Rationale: Windows ships no `sleep`. `timeout /t N` is integer-only and
 * fails with "Input redirection is not supported" when stdin is not a
 * console, which is exactly the case when invoked from goal's `run` (or any
 * other process that pipes stdin). This provides a real executable on %PATH%
 * so `run"sleep""0.25"` works.
 *
 * Build: see ./build.sh (or ./Makefile).
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef _WIN32
#include <windows.h>
#else
#include <errno.h>
#include <time.h>
#endif

static void usage(void) {
    fputs("usage: sleep NUMBER[smhd]...\n", stderr);
}

/* Sleep for `seconds`, chunked so long waits cannot overflow a DWORD of ms. */
static void sleep_seconds(double seconds) {
    while (seconds > 0.0) {
        double chunk = seconds > 1000.0 ? 1000.0 : seconds;
#ifdef _WIN32
        Sleep((DWORD)(chunk * 1000.0 + 0.5));
#else
        struct timespec ts;
        ts.tv_sec = (time_t)chunk;
        ts.tv_nsec = (long)((chunk - (double)ts.tv_sec) * 1e9);
        while (nanosleep(&ts, &ts) == -1 && errno == EINTR) {
            /* resume remaining time */
        }
#endif
        seconds -= chunk;
    }
}

/* Parse one interval argument. Returns 0 on success, -1 on bad input. */
static int parse_interval(const char *arg, double *out) {
    char *end = NULL;
    double v = strtod(arg, &end);

    if (end == arg) {
        return -1;
    }
    switch (*end) {
    case '\0':                        break;
    case 's': end++;                  break;
    case 'm': v *= 60.0;    end++;    break;
    case 'h': v *= 3600.0;  end++;    break;
    case 'd': v *= 86400.0; end++;    break;
    default:
        return -1;
    }
    if (*end != '\0' || v < 0.0) {
        return -1;
    }
    *out = v;
    return 0;
}

int main(int argc, char **argv) {
    double total = 0.0;
    int i;

    if (argc < 2) {
        usage();
        return 2;
    }

    for (i = 1; i < argc; i++) {
        double v;

        if (strcmp(argv[i], "--help") == 0 || strcmp(argv[i], "-h") == 0) {
            usage();
            return 0;
        }
        if (parse_interval(argv[i], &v) != 0) {
            fprintf(stderr, "sleep: invalid time interval '%s'\n", argv[i]);
            return 1;
        }
        total += v;
    }

    sleep_seconds(total);
    return 0;
}
