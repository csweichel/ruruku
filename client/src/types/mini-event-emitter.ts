
export class MiniEventEmitter<T> {
    protected value: T | undefined;
    protected subscriber: (t: T) => void

    public subscribe(sub: (t: T) => void) {
        this.subscriber = sub;
        if (this.value !== undefined) {
            sub(this.value);
        }
    }

    public publish(t: T) {
        this.value = t;
        if (this.subscriber !== undefined) {
            this.subscriber(t);
        }
    }

}
