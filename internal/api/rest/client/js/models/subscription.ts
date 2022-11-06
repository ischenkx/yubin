export default class Subscription {
    public static fromRaw(data: any) {
        data ||= {}
        return new Subscription(
            data['subscriber'] || '',
            data['topic'] || '',
            data['at'] || '',
            data['meta'] || null,
        )
    }

    constructor(
        public subscriber: string,
        public topic: string,
        public at: string,
        public meta: any
    ) {
    }
}