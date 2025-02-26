import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const isValidIP = (ip: string) => {
    // IPv4 validation regex
    const ipv4Pattern = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/;
    if (!ipv4Pattern.test(ip)) return false;
    
    return ip.split('.').map(Number).every(num => num >= 0 && num <= 255);
  };
  
  const isValidCIDR = (cidr: string) => {
    const parts = cidr.split('/');
    if (parts.length !== 2) return false;
    
    const ip = parts[0];
    const prefix = parseInt(parts[1], 10);
    
    return isValidIP(ip) && !isNaN(prefix) && prefix >= 0 && prefix <= 32;
  };
  
  const isValidDomain = (domain: string) => {
    const domainPattern = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
    return domainPattern.test(domain);
  };

export {isValidCIDR, isValidDomain, isValidIP}