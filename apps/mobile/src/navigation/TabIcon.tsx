import React from 'react';
import Svg, { Path, Rect } from 'react-native-svg';

interface Props {
  name: string;
  color: string;
  size?: number;
}

// ไอคอนเส้น (อิง prototype) — ใช้ react-native-svg
export const TabIcon: React.FC<Props> = ({ name, color, size = 23 }) => {
  const common = {
    width: size,
    height: size,
    viewBox: '0 0 24 24',
    fill: 'none',
    stroke: color,
    strokeWidth: 1.7,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
  };

  switch (name) {
    case 'Dashboard':
      return (
        <Svg {...common}>
          <Path d="M4 11l8-7 8 7" />
          <Path d="M6 10v9h12v-9" />
        </Svg>
      );
    case 'Reports':
      return (
        <Svg {...common}>
          <Path d="M5 19V11M12 19V5M19 19v-6" />
        </Svg>
      );
    case 'AI':
      return (
        <Svg {...common}>
          <Path d="M12 3l1.8 5.4L19 10l-5.2 1.6L12 17l-1.8-5.4L5 10l5.2-1.6z" />
        </Svg>
      );
    case 'Hub':
      return (
        <Svg {...common}>
          <Rect x="4" y="4" width="7" height="7" rx="1.6" />
          <Rect x="13" y="4" width="7" height="7" rx="1.6" />
          <Rect x="4" y="13" width="7" height="7" rx="1.6" />
          <Rect x="13" y="13" width="7" height="7" rx="1.6" />
        </Svg>
      );
    default:
      return null;
  }
};
