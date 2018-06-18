import React from "react";
import { Permissions } from "expo";

import TabPage from "./src/Main";

import {
	Database,
	Errors,
	Credentials,
	Interfaces,
} from "keepass.io";

export default class App extends React.Component {
	state = {
		hasFilePermission
		db: new Database()
	}

	constructor() {

	}

	async componentWillMount() {
		const { status } = await Permissions.askAsync(Permissions.FILE);
		this.setState({ hasFilePermission: status === "granted" });
	}

	render() {
		return (<Main />);
	}
}