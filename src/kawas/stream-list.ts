type StreamListSubscriber<T> = {
        push: StreamListPushHandler<T>,
        end: StreamListEndHandler,
};

type StreamListPushHandler<T> = (item: T, is_current: boolean) => void;
type StreamListEndHandler = (is_current: boolean) => void;

export class StreamList<a> implements AsyncIterable<a> {
        backing: a[] = [];
        push_sub: StreamListPushHandler<a> | undefined = undefined;
        end_sub: StreamListEndHandler | undefined = undefined;
        ended = false;
        push(a: a): void {
                if (this.ended) {
                        throw new Error("cannot push to an ended stream");
                }
                if (this.push_sub) {
                        this.push_sub(a, true);
                } else {
                        this.backing.push(a);
                }
        }
        end(): void {
                if (this.ended) {
                        return;
                }
                this.ended = true;
                if (this.end_sub) {
                        this.end_sub(true);
                }
        }
        getBacking(): a[] {
                return this.backing;
        }
        forEach(push_handler: StreamListPushHandler<a>): void {
                if (this.push_sub) {
                        throw new Error("body has been consumed");
                }
                for (const item of this.backing) {
                        push_handler(item, false);
                }
                this.backing.length = 0;  // drop all items in array
                this.push_sub = push_handler;
        }
        onEnd(end_handler: StreamListEndHandler): void {
                if (this.end_sub) {
                        throw new Error("end event already has listener");
                }
                if (this.ended) {
                        end_handler(false);
                }
                this.end_sub = end_handler;
        }
        
        [Symbol.asyncIterator](): AsyncIterator<a, void, undefined> {
                let messageQueue: a[] = [];
                let currentPromise: Promise<IteratorResult<a, void>> | null = null;
                let currentResolve: ((value: IteratorResult<a, void>) => void) | null = null;
                let currentReject: ((reason?: any) => void) | null = null;
                let iteratorFinished = false;

                const iteratorPushHandler: StreamListPushHandler<a> = (item, _is_current) => {
                        messageQueue.push(item);
                        if (currentResolve) {
                                const resolve = currentResolve;
                                currentPromise = null;
                                currentResolve = null;
                                currentReject = null;
                                resolve({ value: messageQueue.shift()!, done: false });
                        }
                };

                const iteratorEndHandler: StreamListEndHandler = () => {
                        iteratorFinished = true;
                        if (currentResolve) {
                                const resolve = currentResolve;
                                currentPromise = null;
                                currentResolve = null;
                                currentReject = null;
                                resolve({ value: undefined, done: true });
                        }
                };

                // subscribe will immediately call iteratorPushHandler for all backedup items
                this.forEach(iteratorPushHandler);
                this.onEnd(iteratorEndHandler);

                return {
                        async next(): Promise<IteratorResult<a, void>> {
                                if (iteratorFinished && messageQueue.length === 0) {
                                        return { value: undefined, done: true };
                                }

                                if (messageQueue.length > 0) {
                                        return { value: messageQueue.shift()!, done: false };
                                }

                                if (this.ended) {
                                        iteratorFinished = true;
                                        return { value: undefined, done: true };
                                }

                                if (!currentPromise) {
                                        currentPromise = new Promise((resolve, reject) => {
                                                currentResolve = resolve;
                                                currentReject = reject;
                                        });
                                }
                                return currentPromise;
                        },

                        async return(): Promise<IteratorResult<a, void>> {
                                iteratorFinished = true;
                                if (currentResolve) {
                                        currentResolve({ value: undefined, done: true });
                                        currentPromise = null;
                                        currentResolve = null;
                                        currentReject = null;
                                }
                                // In a real system, you'd want an `unsubscribe` mechanism here.
                                return { value: undefined, done: true };
                        },

                        async throw(e?: any): Promise<IteratorResult<a, void>> {
                                iteratorFinished = true;
                                if (currentReject) {
                                        currentReject(e);
                                } else if (currentResolve) {
                                        currentResolve({ value: undefined, done: true });
                                }
                                currentPromise = null;
                                currentResolve = null;
                                currentReject = null;
                                return Promise.reject(e);
                        }
                };
        }
}