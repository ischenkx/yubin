import Base from "./base";
import {Axios} from "axios";
import Template from "../models/template";

export default class Templates extends Base {
    constructor(client: Axios) {
        super(client);
    }

    async getTemplates(offset = null, limit = null) {
        return await this.get('/templates', {
            params: {
                offset, limit
            }
        }).then(data => (data || []).map(Template.fromRaw))
    }

    async getTemplate(id: string) {
        return await this.get(`/templates/${id}`).then(Template.fromRaw)
    }

    async deleteTemplate(id: string) {
        return await this.delete(`/templates/${id}`)
    }

    async createTemplate(template: Template) {
        return await this.post(`/templates`, {
            name: template.name,
            data: template.data,
            sub: template.sub,
            meta: template.meta
        })
    }

    async updateTemplate(info: any) {
        return await this.put(`/templates`, info)
    }
}