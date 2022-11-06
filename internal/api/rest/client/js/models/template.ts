export default class Template {
    public static fromRaw(data: any) {
        data ||= {}
        return new Template(
            data['name'] || '',
            data['data'] || '',
            data['meta'] || null,
            data['sub'] || null
        )
    }

    constructor(
        public name: string,
        public data: string,
        public meta: any,
        public sub: any
    ) {}
}