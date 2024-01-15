// TODO httponly?
export const CookieStore = {
    getCookie: function(name: string): string {
        // Read cookie back
        const cookieName = `${name}=`;
        const cookieArray = document.cookie.split(';');
        for (let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i].trim();
            if (cookie.indexOf(cookieName) === 0) {
                return cookie.substring(cookieName.length, cookie.length);
            }
        }
        return "";
    },

    setCookie: function(name: string, value: string, days: number) {
        const expires = new Date();
        expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
        document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;samesite=strict`;
    }
}