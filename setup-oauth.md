# Google OAuth Setup Guide

## Quick Setup (For Testing)

If you want to test the application without setting up OAuth, it will work perfectly with local authentication only. The Google login button will be hidden and users can register/login normally.

## Full OAuth Setup

To enable Google OAuth login, follow these steps:

### 1. Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the **Google+ API** (or People API)

### 2. Create OAuth 2.0 Credentials

1. Go to **APIs & Services** → **Credentials**
2. Click **+ CREATE CREDENTIALS** → **OAuth 2.0 Client IDs**
3. Choose **Web application**
4. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
5. Copy the **Client ID** and **Client Secret**

### 3. Set Environment Variables

#### Option A: Export in terminal (temporary)
```bash
export GOOGLE_CLIENT_ID="your-client-id-here.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="your-client-secret-here"
```

#### Option B: Create .env file (recommended)
```bash
# Create .env file
cp .env.example .env

# Edit .env file and add your credentials
nano .env
```

Then add:
```env
GOOGLE_CLIENT_ID=your-client-id-here.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret-here
```

#### Option C: Load from .env file in Go (optional)
If you want to automatically load from .env file, install and use the godotenv package:

```bash
go get github.com/joho/godotenv
```

### 4. Start the Server

```bash
# With environment variables set
go run simple_server.go
```

You should see:
```
✅ Google OAuth configured successfully
```

Instead of:
```
⚠️  Google OAuth not configured
```

### 5. Test OAuth Login

1. Go to `http://localhost:8080/login`
2. You should now see the "Login with Google" button
3. Click it to test the OAuth flow

## Security Notes

- Keep your Client Secret secure and never commit it to version control
- For production, use HTTPS and update the redirect URI accordingly
- Consider using environment-specific configuration files
- Rotate credentials periodically

## Troubleshooting

### "OAuth not configured" message
- Verify environment variables are set: `echo $GOOGLE_CLIENT_ID`
- Restart the server after setting environment variables

### "Authorization Error: Missing client_id"
- Client ID environment variable is not set or empty
- Check for typos in environment variable names

### "Redirect URI mismatch"
- Verify the redirect URI in Google Console matches exactly: `http://localhost:8080/auth/google/callback`
- Make sure there are no trailing slashes or extra characters

### "API not enabled" 
- Enable the Google+ API or People API in Google Cloud Console
- Wait a few minutes for the API to become available

## Testing Without OAuth

The application works perfectly without OAuth! Users can:
- Register new accounts with email/password
- Login with username/password  
- Use all workout tracking features
- Access admin panel (if admin role)

Only the Google login option will be hidden when OAuth is not configured.
