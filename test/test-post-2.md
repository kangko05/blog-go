# React TypeScript 프로젝트 설정

React와 TypeScript를 함께 사용하여 타입 안전한 웹 애플리케이션을 만드는 방법을 정리했습니다.

## 프로젝트 생성

```bash
npx create-react-app my-app --template typescript
cd my-app
npm start
```

## 주요 설정 파일

### tsconfig.json

TypeScript 컴파일러 설정을 조정합니다:

```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true
  }
}
```

## 컴포넌트 타입 정의

```typescript
interface Props {
  title: string;
  count: number;
}

const MyComponent: React.FC<Props> = ({ title, count }) => {
  return (
    <div>
      <h1>{title}</h1>
      <p>Count: {count}</p>
    </div>
  );
};
```

## Hooks와 TypeScript

useState, useEffect 등의 Hooks를 TypeScript와 함께 사용할 때의 패턴들을 익히는 것이 중요합니다.
