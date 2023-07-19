export default class PersonalReport {
    public static fromRaw(data: any) {
        data ||= {}
        return new PersonalReport(
            data['publication_id'] || '',
            data['user_id'] || '',
            data['status'] || '',
            data['meta'] || null
        )
    }

    constructor(
        public publicationId: string,
        public userId: string,
        public status: string,
        public meta: any
    ) {
    }
}