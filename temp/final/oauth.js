const CLIENT_ID = '194157837347-321ivntptlcfssjhveeeci6rt03ocsa0.apps.googleusercontent.com';
const API_KEY = 'AIzaSyCLnLvZLgtKGoGcC6fEP1euO3KiXl2YWOY';
const DISCOVERY_DOC = 'https://www.googleapis.com/discovery/v1/apis/calendar/v3/rest';
const SCOPES = 'https://www.googleapis.com/auth/calendar.readonly';

let tokenClient;

function gapiLoaded() {
	gapi.load('client', initializeGapiClient);
}

async function initializeGapiClient() {
	await gapi.client.init({
		apiKey: API_KEY,
		discoveryDocs: [DISCOVERY_DOC],
	});
}

function gisLoaded() {
	tokenClient = google.accounts.oauth2.initTokenClient({
		client_id: CLIENT_ID,
		scope: SCOPES,
		callback: '', // defined later
	});
}

function handleAuthClick() {
	tokenClient.callback = async (resp) => {
		if (resp.error !== undefined) {
			throw (resp);
		}
		await gen_schedule();
	};

	if (gapi.client.getToken() === null) {
		tokenClient.requestAccessToken({prompt: 'consent'});
	} else {
		tokenClient.requestAccessToken({prompt: ''});
	}
}
