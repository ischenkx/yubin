import Base from "./base";
import {Axios} from "axios";
import Source from "../models/source";

export default class Sources extends Base {
    constructor(client: Axios) {
        super(client);
    }

    async getSources(offset = null, limit = null) {
        return await this
            .get('/sources', {
                params: {
                    offset, limit
                }
            })
            .then(data => (data || []).map(Source.fromRaw))
    }

    async getSource(name: string) {
        return await this.get(`/sources/${name}`).then(Source.fromRaw)
    }

    async deleteSource(name: string) {
        return await this.delete(`/sources/${name}`)
    }

    async createSource(source: Source) {
        return await this
            .post('/sources', {
                name: source.name,
                address: source.address,
                password: source.password,
                host: source.host,
                port: source.port
            })
            .then(Source.fromRaw)
    }

    async updateSource(info: any) {
        return await this
            .put('/sources', info)
            .then(Source.fromRaw)
    }
}