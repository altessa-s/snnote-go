// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

// The C side of SetLogger: a sink with the signature the library calls, which
// hands each line to the exported Go function. Kept apart from log.go because
// a Go file that exports to C may only declare, not define, in its preamble.

#include "_cgo_export.h"
#include "snnote/log.h"
#ifdef SNNOTE_GO_TEST
#include <stdatomic.h>
static _Atomic int sngo_installs;
int sngo_log_installs(void) { return atomic_load(&sngo_installs); }
#endif

static void sngo_log_sink(sn_log_level level, const char *message, void *context) {
    (void)context;
    sngoLog((int)level, (char *)message);
}

void sngo_log_install(void) {
#ifdef SNNOTE_GO_TEST
    atomic_fetch_add(&sngo_installs, 1);
#endif
    sn_log_install(sngo_log_sink, nullptr);
}
