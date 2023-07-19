import Sources from "../controllers/sources";

export default class Source {
    public static fromRaw(data: any) {
        data ||= {}
        return new Source(
            data['name'] || '',
            data['address'] || '',
            data['password'] || '',
            data['host'] || '',
            data['port'] || 0
        )
    }

    constructor(
        public name: string,
        public address: string,
        public password: string,
        public host: string,
        public port: number
    ) {
    }
}
