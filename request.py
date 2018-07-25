import argparse
import requests

def requestWeb(url):
    http_proxy  = "http://127.0.0.1:8000"
    https_proxy = "https://127.0.0.1:8000"
    ftp_proxy   = "ftp://127.0.0.1:8000"

    proxyDict = { 
                "http"  : http_proxy, 
                "https" : https_proxy, 
                "ftp"   : ftp_proxy
                }

    r = requests.get(url, proxies=proxyDict)
    print(r.text)

def main():
    args = load_input()
    requestWeb(args.i)

def load_input():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        '-i', '--ip', help='IP of the request website', required=True, type=str)
    return parser.parse_args()

if __name__ == '__main__':
    main()