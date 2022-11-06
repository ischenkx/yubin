export default class Report {
    public static fromRaw(data: any) {
        data ||= {}
        return new Report(
            data['publication_id'] || '',
            data['status'] || '',
            data['failed'] || [],
            data['ok'] || []
        )
}
    constructor(
        public publicationId: string,
        public status: string,
        public failed: [string],
        public ok: [string]
    ) {}
}