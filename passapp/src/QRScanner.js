import React from 'react'

import {
	Text,
	StyleSheet,
	BackHandler,
} from 'react-native'

import { 
	BarCodeScanner,
	Permissions,
} from 'expo';

import { Buffer } from "buffer";
import nacl from "tweetnacl";

export default class QRScanner extends React.Component {
	state = {
		hasCameraPermission: null,
	}

	constructor(props) {
		super(props)

		BackHandler.addEventListener('hardwareBackPress', this.exit)
	}

	exit = () => {
		BackHandler.removeEventListener('hardwareBackPress', this.exit)
		this.props.onExit()
		return true
	}

	async componentWillMount() {
		const { status } = await Permissions.askAsync(Permissions.CAMERA);
		this.setState({ hasCameraPermission: status === "granted" });
	}

	_callback = ({ type, data }) => {
		bytes = Buffer.from(data, "base64");
		if (bytes.length != 128) return;

		clientSign = bytes.slice(64, 96);
		clientBox = bytes.slice(96, 128);

		verify = nacl.sign.open(
			bytes,
			clientSign,
		);

		if (verify === null) return;
		this.props.onFind({
			clientSign: clientSign, 
			clientBox: clientBox,
		});
	}

	render() {
		if (this.state.hasCameraPermission === null) {
			return <Text>Requesting for camera permission</Text>;
		} else if (this.state.hasCameraPermission === false) {
			this.exit();
			return;
		} else {
			return (
				<BarCodeScanner
					onBarCodeRead={this._callback}
					style={StyleSheet.absoluteFill}
				/>
			);
		}
	}
}