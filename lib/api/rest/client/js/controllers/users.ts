import {Axios} from "axios";
import Base from "./base";
import User from "../models/user";

export default class Users extends Base {
    constructor(client: Axios) {
        super(client)
    }

    async getUsers(offset : number | null = null, limit : number | null = null) {
        return await this
            .get('/users', {
                params: {
                    offset, limit
                }
            })
            .then(data => (data || []).map(User.fromRaw))
    }

    async getUser(id: string) {
        return await this
            .get(`/users/${id}`)
            .then(User.fromRaw)
    }

    async deleteUser(id: string) {
        return await this
            .delete(`/users/${id}`)
    }

    async updateUser(info: any) {
        return await this
            .put(`/users`, info)
            .then(User.fromRaw)
    }

    async createUser(user: User) {
        return await this
            .post(`/users`, {
                name: user.name,
                surname: user.surname,
                email: user.email
            })
            .then(User.fromRaw)
    }
}