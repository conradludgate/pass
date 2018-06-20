import React from "react";
import { StyleSheet, Text, View, Button } from "react-native";

import { BarCodeScanner, Permissions } from "expo";

import { Buffer } from "buffer";
import nacl from "tweetnacl";

import QRScanner from "./QRScanner"

export default class TabPage extends React.Component {
	state = {
		scan: false,
	};

	_onPress = () => {
		this.setState({scan:true})
	}

	_onFind = ({ clientSign, clientBox }) => {
		this.setState({scan:false});

		alert("sign: " + clientSign.toString("base64") + ", box: " + clientBox.toString("base64"))
	};

	_onExit = (error) => {
		this.setState({scan:false});

		if (!!error) alert(error);
	}

	render() {
		if (this.props.page === "1") {
			if (this.state.scan) {
				return (<QRScanner
					onFind={this._onFind}
					onExit={this._onExit}
				/>);
			} else {
				return (<Button 
					onPress={this._onPress}
					title="Scan"
				/>);
			}

		} else if (this.props.page === "2") {
			return <Text>Passwords Page</Text>;
		}
	}
}