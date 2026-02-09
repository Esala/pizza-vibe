import { NextRequest, NextResponse } from 'next/server';

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ item: string }> }
) {
  try {
    const { item } = await params;
    const inventoryServiceUrl = process.env.INVENTORY_SERVICE_URL || 'http://localhost:8084';
    const response = await fetch(`${inventoryServiceUrl}/inventory/${item}`, {
      method: 'POST',
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to acquire item' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error acquiring item:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
