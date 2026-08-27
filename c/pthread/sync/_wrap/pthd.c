#include <pthread.h>

// -----------------------------------------------------------------------------

pthread_once_t llgoSyncOnceInitVal = PTHREAD_ONCE_INIT;

typedef void (*llgo_once_context_fn)(void *);

typedef struct llgo_once_context_frame {
    llgo_once_context_fn callback;
    void *data;
    struct llgo_once_context_frame *previous;
} llgo_once_context_frame;

static _Thread_local llgo_once_context_frame *llgo_once_context_top;

static void llgo_once_context_invoke(void) {
    llgo_once_context_frame *frame = llgo_once_context_top;
    llgo_once_context_top = frame->previous;
    frame->callback(frame->data);
}

int llgo_pthread_once_context(pthread_once_t *once_control,
                              llgo_once_context_fn callback, void *data) {
    llgo_once_context_frame frame = {
        callback,
        data,
        llgo_once_context_top,
    };
    int result;

    llgo_once_context_top = &frame;
    result = pthread_once(once_control, llgo_once_context_invoke);
    if (llgo_once_context_top == &frame)
        llgo_once_context_top = frame.previous;
    return result;
}

// -----------------------------------------------------------------------------

// wrap return type to void
void wrap_pthread_mutex_lock(pthread_mutex_t *mutex) {
    pthread_mutex_lock(mutex);
}

// wrap return type to void
void wrap_pthread_mutex_unlock(pthread_mutex_t *mutex) {
    pthread_mutex_unlock(mutex);
}

// -----------------------------------------------------------------------------
