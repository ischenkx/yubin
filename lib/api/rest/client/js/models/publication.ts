export default class Publication {
    public static fromRaw(data: any) {
        data ||= {}
        return new Publication(
            data["source"] || '',
            data["template"] || '',
            data["users"] || [],
            data['topics'] || [],
            data["meta"] || null,
            data["at"] || 0
        )
    }

    constructor(
        public source: string,
        public template: string,
        public users: string[] | null,
        public topics: string[] | null,
        public meta: any,
        public at: number) {
    }
}