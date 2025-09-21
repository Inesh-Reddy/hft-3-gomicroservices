import { Module } from '@nestjs/common';
import { TickerController } from './ticker.controller';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'TICKER_PACKAGE',
        transport: Transport.GRPC,
        options: {
          package: 'ticker',
          protoPath: join(__dirname, '../../../../packages/proto/ticker.proto'),
          url: 'localhost:50052',
        },
      },
    ]),
  ],
  controllers: [TickerController],
})
export class TickerModule {}
